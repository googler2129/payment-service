# Payment Service - SCD Implementation

## Introduction

This service implements a payment processing system with a focus on handling Slowly Changing Dimension (SCD) Type 2 data. The solution in the `pkg/repository/scd` folder provides a generic abstraction layer that entirely hides the complexity of SCD from developers.

## Understanding the Problem

### What is SCD Type 2?

SCD Type 2 is a data warehousing pattern used to track historical changes to records. In our implementation:

- Each record has a unique `id` that remains constant across versions
- Each version of the record gets a unique `version` number
- Each record version also has a unique `uid` to serve as a primary key
- The `is_latest` flag identifies the current version of a record

For example, in our Job table, changes to status or rate create new versions while preserving history:

| id | version | uid | status | rate | title |
| --- | --- | --- | --- | --- | --- |
| job_123 | 1 | job_uid_abc | extended | 20.0 | Software Engineer |
| job_123 | 2 | job_uid_def | active | 20.0 | Software Engineer |
| job_123 | 3 | job_uid_ghi | active | 15.5 | Software Engineer |

### The Challenges We Faced

Without proper abstraction, SCD tables create several challenges:

1. **Complex Queries**: Developers must always filter for the latest version or implement complex joins
2. **Code Repetition**: The same SCD handling logic gets duplicated across repositories
3. **Performance Concerns**: Inefficient SCD queries can be slow, especially with millions of records
4. **Inconsistent Implementation**: Different teams might handle versions differently

These challenges were evident in queries like:
- Get all active Jobs for a company (latest version only)
- Get all PaymentLineItems for a contractor in a time period (latest versions only)
- Get all Timelogs for a particular job (latest versions only)

## Our Solution

### Why We Chose a Generic SCD Repository

We implemented a generic SCD repository pattern leveraging Go's generics, producing multiple benefits:

1. **Generic Implementation**: One implementation works for any entity requiring versioning (Jobs, Timelogs, PaymentLineItems)
2. **Type Safety**: The generic approach maintains type safety while providing SCD functionality
3. **Unified Interface**: All SCD operations follow the same pattern regardless of entity type
4. **Performance Optimization**: The `is_latest` column enables efficient querying
5. **Proper Abstraction**: SCD complexity is hidden from application code

### How Our Implementation Solves the Problem

#### Zero Code in Repositories

After implementing our SCD abstraction, developers **don't need to write any code in particular repositories** to handle SCD functionality. The repositories become extremely thin wrappers that simply embed the generic implementation, resulting in minimal boilerplate and consistent behavior across all versioned entities.

**Important**: After using the SCD generic repository, repositories no longer need to write any code for common CRUD operations (GET, CREATE, DELETE, UPDATE). All these operations are inherited from the generic implementation, eliminating redundancy and ensuring consistent behavior.

#### The Job Repository Example

Our solution completely abstracts SCD handling as shown in the job repository implementation:

```go
type JobRepository struct {
    scd.SCDRepository[domain.Job]
    db *postgres.DbCluster
}

func NewJobRepository(db *postgres.DbCluster) domain.JobRepositoryInterface {
    repoOnce.Do(func() {
        repo = &JobRepository{
            db:            db,
            SCDRepository: scd.NewSCDRepository(db, domain.Job{}),
        }
    })

    return repo
}
```

By embedding `scd.SCDRepository[domain.Job]`, the repository instantly inherits all SCD functionality with zero boilerplate code. This pattern is consistent across all SCD-enabled entities.

#### Eliminated Code Redundancy

Without our generic solution, developers would need to implement these patterns repeatedly for each entity:

```go
// Without our abstraction (example of what we ELIMINATED):
func (r *jobRepository) FindLatestJobs(ctx context.Context, filter map[string]interface{}) ([]Job, error) {
    var results []Job
    // Complex subquery to find latest version
    err := r.db.GetSlaveDB(ctx).Raw(`
        SELECT j.* FROM job j
        INNER JOIN (
            SELECT id, MAX(version) as max_version 
            FROM job
            GROUP BY id
        ) latest ON j.id = latest.id AND j.version = latest.max_version
        WHERE /* custom filters */
    `).Find(&results).Error
    // Error handling...
    return results, nil
}
```

#### Developer Experience

The SCD repository provides a clean, intuitive interface that abstracts away all versioning complexity:

```go
// Creating a new record (version 1 automatically assigned)
repo.Create(ctx, jobRecord)

// Updating a record (new version created with incremented version number)
repo.Update(ctx, jobId, updatedJobRecord)

// Finding only the latest version - no SCD complexity visible
job, err := repo.FindByID(ctx, jobId)

// Finding records with filters (only latest versions)
jobs, err := repo.FindLatestWithFilter(ctx, map[string]interface{}{
    "status": "active",
    "company_id": companyId,
})

// For historical analysis, access all versions
jobVersions, err := repo.FindVersionsForID(ctx, jobId)
```

### Performance Optimization with `is_latest` Flag

We deliberately added an `is_latest` column as a performance optimization for our read-heavy workload:

1. **Optimized for Reads**: For read-heavy workloads, the `is_latest` flag allows direct indexing and efficient filtering
2. **Fast Latest Record Retrieval**: Retrieving latest versions requires a simple WHERE clause instead of complex subqueries
3. **Simpler Queries**: Application code can use simple conditions without complex JOINs and subqueries

Example query transformation:

**Before our abstraction**:
```sql
-- Complex query with subqueries
SELECT j.* FROM job j
INNER JOIN (
    SELECT id, MAX(version) as max_version 
    FROM job
    GROUP BY id
) latest ON j.id = latest.id AND j.version = latest.max_version
WHERE status = 'active' AND company_id = 'comp_123';
```

**With our abstraction**:
```sql
-- Simple, efficient query
SELECT * FROM job 
WHERE is_latest = true AND status = 'active' AND company_id = 'comp_123';
```

For high-write/low-read scenarios, the alternative approach would be more optimal, but our service prioritizes read performance.

### Generic Static Repository

For entities that don't require versioning (like Contractor), we've created a parallel generic static repository pattern:

```go
type ContractorRepository interface {
    static.StaticRepository[Contractor]
}
```

#### Static Repository Example

**Before (Old approach - writing custom methods for each entity):**
```go
// Before: Each repository had to implement its own GET methods
type ContractorRepository interface {
    Create(ctx context.Context, contractor *Contractor) error
    Update(ctx context.Context, contractor *Contractor) error
    Delete(ctx context.Context, contractor *Contractor) error
    GetByID(ctx context.Context, id string) (*Contractor, error)
    GetByCompanyID(ctx context.Context, companyID string) ([]Contractor, error)
    // Many more custom methods...
}

// Implementation required writing each method
type contractorRepositoryImpl struct {
    db *postgres.DbCluster
}

func (r *contractorRepositoryImpl) GetByID(ctx context.Context, id string) (*Contractor, error) {
    var result Contractor
    err := r.db.GetSlaveDB(ctx).Where("id = ?", id).First(&result).Error
    return &result, err
}

func (r *contractorRepositoryImpl) GetByCompanyID(ctx context.Context, companyID string) ([]Contractor, error) {
    var results []Contractor
    err := r.db.GetSlaveDB(ctx).Where("company_id = ?", companyID).Find(&results).Error
    return results, err
}

// And so on for each method...
```

**After (New approach - using generics):**
```go
// After: Simply embed the generic interface
type ContractorRepository interface {
    static.StaticRepository[Contractor]
}

// Contractor must implement the Static interface
type Contractor struct {
    ID        string `gorm:"column:id;primaryKey"`
    Name      string `gorm:"column:name"`
    CompanyID string `gorm:"column:company_id"`
    // Other fields...
}

// Implement the SetID method required by the Static interface
func (c *Contractor) SetID(id string) {
    c.ID = id
}

// Creating the repository is now a one-liner
func NewContractorRepository(db *postgres.DbCluster) ContractorRepository {
    return static.NewStaticRepository[Contractor](db)
}

// Using the repository
func ExampleUsage(ctx context.Context, repo ContractorRepository) {
    // Create a contractor
    contractor := &Contractor{
        Name:      "John Doe",
        CompanyID: "comp_123",
    }
    repo.Create(ctx, contractor)
    
    // Get by conditions
    filter := map[string]interface{}{"company_id": "comp_123"}
    contractors, err := repo.GetAllByConditions(ctx, filter)
    
    // Get single record
    singleFilter := map[string]interface{}{"id": "some-id"}
    contractorRecord, err := repo.GetByConditions(ctx, singleFilter)
    
    // Update
    if contractorRecord != nil {  // Note: checking pointer directly, not dereferencing
        repo.UpdateByCondition(ctx, 
            map[string]interface{}{"id": (*contractorRecord).ID},
            contractorRecord,
        )
    }
}
```

This removes all the boilerplate code in repositories and no code would be needed for GET/UPDATE/CREATE/DELETE operations. 

### Cross-Language Compatibility Solution

To address the challenge of supporting both Django and Golang services, we can use database VIEWS:

```sql
CREATE VIEW latest_jobs AS
SELECT * FROM job WHERE is_latest = true;
```


## Getting Started

### Running the PostgreSQL Database

```bash
docker compose up
```

### Running Migrations

```bash
CONFIG_SOURCE=local go run main.go --mode=migration
```

### Running the Service

```bash
CONFIG_SOURCE=local go run main.go
```