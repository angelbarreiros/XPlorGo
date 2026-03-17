# XPlorGo SDK

## Overview

**XPlorGo** is a Go SDK to interact with the **XPlor API** from Resamania. It provides a partial interface to manage resources in a system for events, activities, classes, and users in sports clubs.

---

## Initialization and Authentication

```go
provider := xplorcore.Init(xplorcore.NewConfig(
    host string,              // "host
    version string,           // "endpoint"
    enterpriseName string,    // "org_name"
    clientId string,          // OAuth2 Client ID
    clientSecret string,      // OAuth2 Client Secret
    headers map[string]string // Additional headers (e.g. API keys)
))

provider.Close() // Releases resources
```

- Automatic OAuth2 authentication
- Reuses tokens while they are valid
- Thread-safe with automatic synchronization

---

## Main Entities

### 1. **Contacts** (System Customers/Users)
```
Entity: XPlorContact
- ContactID (unique ID)
- Email, FamilyName, GivenName
- BirthDate, Gender
- Address (full address)
- ClubID (associated club)
- Mobile, NationalID
- State (contact status)
- PictureID, SourceID
- SalepersonID (associated salesperson)
- CreatedAt, UpdatedAt

Useful methods:
- ContactID() -> string
- FullAddress() -> string (full address)
- Age() -> int (calculates age from birthDate)
- ClubIDValue() -> string
- InitialSalepersonIDValue() -> string
```

**Search parameters:**
```
XPlorContactsParams:
- ContactID / ContactIDs[]
- ClubID / ClubIDs[]
- State / States[]
- Email / Emails[]
- Mobile
- Number
- FamilyName
- GivenName
```

---

### 2. **Activities**
```
Entity: XPlorActivity
- ActivityID
- Name (e.g. "PADEL", "Fitness")
- ClubID
- ColorHex (activity color)
- Durations[] (ISO 8601: PT60M, PT90M)
- IsBookable, IsViewable
- ShowcaseActivities[]
- ActivityGroups[]
- ArchivedAt/By
- CreatedAt/By
- TemplateToken

Useful methods:
- ActivityID() -> string
- ShowcaseIDs() -> []int
- IsActive() -> bool (not archived)
- IsPadel() -> bool
- DurationMinutes() -> int (extracts minutes from PT format)
```

**Search parameters:**
```
XPlorActivitiesParams:
- ClubID / ClubIDs[]
- Name
- Archived (bool)
```

---

### 3. **Classes** (Class Events)
```
Entity: XPlorClass
- ClassEventID
- Club, Studio, Activity, Coach
- StartedAt, EndedAt (LocalTime)
- AttendingLimit, QueueLimit
- BookedAttendees[]
- QueuedAttendees[]
- AttendeeRemaining
- QueueRemaining
- Summary, Description
- CoachAvailable (bool)
- DisabledItems[]
- Recurrence
- ClassLayout, ClassLayoutConfiguration
- Processing (bool)
- CreatedAt/UpdatedAt/DeletedAt/ArchivedAt

Useful methods:
- ClassEventID() -> string
- ClubID() -> string
- StudioID() -> string
- ActivityID() -> string
- CoachID() -> string
- RecurrenceID() -> string
- GetStartedAt() -> time.Time
- GetEndedAt() -> time.Time
- HasAvailableSpots() -> bool
- HasQueueSpots() -> bool
- IsActive() -> bool (not deleted and not archived)
- IsDeleted() -> bool
- GetAllContactIDs() -> []string
```

**Search parameters:**
```
XPlorClassesParams:
- Club / Clubs[]
- Coach / Coaches[]
- Activity / Activities[]
- Studio / Studios[]
- RecurrenceID / RecurrenceIDs[]
- StartAt (LocalDateTime)
- EndAt (LocalDateTime)
- StateCode
```

---

### 4. **Studios**
```
Entity: XPlorStudio
- StudioID
- Name
- Club
- ZoneID
- Capacity
- Overbooking
- StreetAddress, PostalCode
- AddressLocality, AddressCountry
- Tags
- CreatedAt/By
- ArchivedAt/By

Useful methods:
- StudioID() -> string
- ClubID() -> string
- ZoneID() -> string
- Address() -> string
```

---

### 5. **Coaches**
```
Entity: XPloreCoach
- CoachID
- GivenName, FamilyName
- AlternateName
- Email, Mobile
- Activities[]
- CreatedAt/By
- ArchivedAt/By

Useful methods:
- CoachID() -> string
- ActivityIDs() -> []string
```

---

### 6. **Clubs**
```
Entity: XPlorClub
- ClubID, ClubNumberID
- Code (3-5 characters)
- Name
- Number
- Email, Phone
- StreetAddress, PostalCode
- AddressLocality, AddressCountry
- AddressCountryIso
- OpeningDate
- Description
- ClubTags[]
- PublicMetadata
- SaleTerms[]
- Locale
- CreatedAt/By
- DeletedAt

Related fields:
- TaxRates
- ResaboxNotification
```

---

### 7. **Subscriptions**
```
Entity: XPlorSubscription
- SubscriptionID
- Contact
- Club
- StartedAt, EndedAt
- IsActive (bool)
- CreatedAt/UpdatedAt
- [Other plan-specific fields]

Useful methods:
- SubscriptionID() -> string
- ContactID() -> string
- ClubID() -> string
- IsActive() -> bool
- IsExpired() -> bool
```

**Search parameters:**
```
XPlorSubscriptionsParams:
- ContactID / ContactIDs[]
- ClubID
- Active (bool)
- StartedAt / EndedAt (date range)
```

---

### 8. **System Users**
```
Entity: XPlorUser
- UserID
- Email
- GivenName, FamilyName
- Mobile, PictureLink
- Code
- ClubIds[]
- NetworkNodeIds[]
- Roles[]
- Active (bool)
- Locale
- CreatedAt/DeletedAt/ArchivedAt
- [Other fields]

Useful methods:
- UserID() -> string
- ClubIDs() -> []string
- NetworkNodeIDs() -> []string
- PropertiesNetworkNodeIDs() -> []string
- IsActive() -> bool
- IsDeleted() -> bool
- IsArchived() -> bool
- IsInactive() -> bool
- FullName() -> string
- HasRole(role string) -> bool
- GetCreatedAt() -> time.Time
- GetDeletedAt() -> *time.Time
```

---

### 9. **Network Nodes**
```
Entity: XPlorNetworkNode
- NodeID
- Name
- Type
- RelatedClubs[]
- [Other fields]
```

---

### 10. **Events**
```
Entity: XPlorEvent
- EventID
- Name
- StartAt, EndAt
- Location
- [Other fields]
```

---

### 11. **Families**
```
Entity: XPlorFamily
- FamilyID
- Members[]
- CreatedAt/UpdatedAt
- [Other fields]

Search parameters:
XPlorFamiliesParams:
- FamilyID / FamilyIDs[]
- ContactID
```

---

### 12. **Other Entities**

- **Recurrences**: recurrence patterns for classes (daily, weekly, etc.)
- **Class Types**: available class types
- **Counter Lines**: cash/counter lines
- **Contact Tags**: tags to categorize contacts
- **Articles**: system content/articles
- **Zones**: areas inside a club

---

## Public Provider Functions

### Contact Management
```go
Contacts(nodeId string, params *XPlorContactsParams,
         pagination *XPlorPagination) -> (*XPlorContacts, error)
Contact(nodeId string, contactId string) -> (*XPlorContact, error)
```

### Activity Management
```go
Activities(nodeId string, queryParams *XPlorActivitiesParams,
          pagination *XPlorPagination) -> (*XPlorActivities, error)
Activity(nodeId string, activityId string) -> (*XPlorActivity, error)
```

### Class Management
```go
Classes(nodeId string, params *XPlorClassesParams,
       pagination *XPlorPagination) -> (*XPlorClasses, error)
Class(nodeId string, classId string) -> (*XPlorClass, error)
```

### Studio Management
```go
Studios(nodeId string, pagination *XPlorPagination) -> (*XPlorStudios, error)
Studio(nodeId string, studioId string) -> (*XPlorStudio, error)
```

### Coach Management
```go
Coaches(nodeId string, pagination *XPlorPagination) -> (*XPloreCoaches, error)
Coach(nodeId string, coachId string) -> (*XPloreCoach, error)
```

### Club Management
```go
Clubs(nodeId string) -> (*XPloreClubs, error)
Club(nodeId string, clubId string) -> (*XPlorClub, error)
```

### Subscription Management
```go
Subscriptions(nodeId string, params *XPlorSubscriptionsParams,
             pagination *XPlorPagination) -> (*XPlorSubscriptions, error)
Subscription(nodeId string, subscriptionId string) -> (*XPlorSubscription, error)
```

### User Management
```go
Users(pagination *XPlorPagination) -> (*XPlorUsers, error)
User(userId string) -> (*XPlorUser, error)
```

### Network Node Management
```go
NetworkNodes(pagination *XPlorPagination) -> (*XPlorNetworkNodes, error)
NetworkNode(nodeId string) -> (*XPlorNetworkNode, error)
```

### Attendee Management
```go
Attendees(nodeId string, classId *string,
         pagination *XPlorPagination) -> (*XPlorAttendees, error)
```

### Event Management
```go
Events(nodeId string, pagination *XPlorPagination,
      timeGap *XPlorTimeGap) -> (*XPlorEvents, error)
```

### Family Management
```go
Families(nodeId string, params *XPlorFamiliesParams,
        pagination *XPlorPagination) -> (*XPlorFamilies, error)
Family(nodeId string, familyId string) -> (*XPlorFamily, error)
```

### Recurrence Management
```go
Recurrences(nodeId string, params *XPlorRecurrencesParams,
           pagination *XPlorPagination) -> (*XPlorRecurrences, error)
Recurrence(nodeId string, recurrenceId string) -> (*XPlorRecurrence, error)
```

### Contact Tag Management
```go
ContactTags(nodeId string, params *XPlorContactTagsParams,
           pagination *XPlorPagination) -> (*XPlorContactTags, error)
ContactTag(nodeId string, contactTagId string) -> (*XPlorContactTag, error)
```

### Class Type Management
```go
ClassType(nodeId string, classTypeId string) -> (*XPlorClassType, error)
```

### Counter Management
```go
CounterLines(nodeId string, pagination *XPlorPagination) -> (*XPlorCounterLines, error)
CounterLine(nodeId string, counterLineId string) -> (*XPlorCounterLine, error)
```

### Article Management
```go
Articles(nodeId string, pagination *XPlorPagination) -> (*XPlorArticles, error)
Article(nodeId string, articleId string) -> (*XPlorArticle, error)
```

### Zone Management
```go
Zones(nodeId string, params *XPlorZonesParams,
     pagination *XPlorPagination) -> (*XPlorZones, error)
Zone(nodeId string, zoneId string) -> (*XPlorZone, error)
```

### Close
```go
Close() -> void (releases provider resources)
```

---

## Common Parameters

### Pagination
```go
type XPlorPagination struct {
    Page         int // Page number (1-indexed)
    ItemsPerPage int // Items per page
}
```

### Time Range
```go
type XPlorTimeGap struct {
    StartAt *time.Time // Start date/time
    EndAt   *time.Time // End date/time
}
```

### Query Parameter Builders
```go
// For pagination
BuildPaginationQueryParams(pagination) -> url.Values

// For pagination + time range
BuildPaginationAndTimeGapParams(pagination, timeGap) -> url.Values
```

---

## General Usage Pattern

```go
// 1. Initialize the provider
provider := xplorcore.Init(config)

// 2. Call public functions with nodeId
contacts, err := provider.Contacts(nodeId, params, pagination)
if err != nil {
    // Handle error
}

// 3. Work with results
for _, contact := range contacts.Contacts {
    contactID, _ := contact.ContactID()
    fmt.Println(contact.GivenName, contact.FamilyName, contactID)
}

// 4. Close when done
provider.Close()
```

---

## Security Features

1. **OAuth2 Authentication**: secure credentials
2. **Token Caching**: reuses valid tokens
3. **Thread Safety**: automatic synchronization with mutex
4. **Executor Pool**: reuses HTTP connections
5. **Dynamic Headers**: injects user context (Node ID, Club ID)

---

## Implementation Notes

- Functions return `(*Entity, error)` or `(*EntityCollection, error)`
- `nodeId` is required for most functions (identifies the location)
- Some methods do not require `nodeId` (e.g. `Users`, `NetworkNodes`)
- IDs are automatically extracted from IRI paths (e.g. `"/enjoy/clubs/1249" -> "1249"`)
- All dates use `LocalTime` or `LocalDate` with automatic parsing
- The SDK handles automatic synchronization for concurrent calls

---

## Full Example

```go
package main

import (
    "github.com/angelbarreiros/XPlorGo/xplorcore"
    "github.com/angelbarreiros/XPlorGo/xplorentities"
)

func main() {
    // Initialize
    provider := xplorcore.Init(xplorcore.NewConfig(
        "gateway.prod.gravitee.stadline.tech",
        "resa2-mfr",
        "org_name",
        "client_id",
        "client_secret",
        map[string]string{
            "X-gravitee-api-key": "api_key",
            "grant_type": "client_credentials",
        },
    ))

    // Fetch activities with pagination
    activities, err := provider.Activities(
        "2675", // nodeId
        &xplorentities.XPlorActivitiesParams{},
        &xplorentities.XPlorPagination{
            Page:         1,
            ItemsPerPage: 10,
        },
    )

    if err != nil {
        panic(err)
    }

    // Process activities
    for _, activity := range activities.Activities {
        println(activity.Name)
    }

    // Close
    provider.Close()
}
```

---
