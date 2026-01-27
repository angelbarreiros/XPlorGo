# XPlorGo SDK

## Descripción General

**XPlorGo** es un SDK en **Go** para interactuar con la **API XPlor** de Stadline. Proporciona una interfaz completa para gestionar recursos de un sistema de gestión de eventos, actividades, clases y usuarios en clubes deportivos.

---

## Inicialización y Autenticación

```go
provider := xplorcore.Init(xplorcore.NewConfig(
    host string,              // "gateway.prod.gravitee.stadline.tech"
    version string,           // "resa2-mfr"
    enterpriseName string,    // "maisqueauga"
    clientId string,          // OAuth2 Client ID
    clientSecret string,      // OAuth2 Client Secret
    headers map[string]string // Headers adicionales (ej: API keys)
))

provider.Close() // Libera recursos
```

- Autenticación OAuth2 automática
- Reutiliza tokens mientras sean válidos
- Thread-safe con sincronización automática

---

## 📊 Entidades Principales

### 1. **Contactos** (Customers/Usuarios del Sistema)
```
Entidad: XPlorContact
├── ContactID (ID único)
├── Email, FamilyName, GivenName
├── BirthDate, Gender
├── Address (Dirección completa)
├── ClubID (Club asociado)
├── Mobile, NationalID
├── State (estado del contacto)
├── PictureID, SourceID
├── SalepersonID (vendedor asociado)
└── CreatedAt, UpdatedAt

Métodos útiles:
- ContactID() → string
- FullAddress() → string (dirección completa)
- Age() → int (calcula edad desde birthDate)
- ClubIDValue() → string
- InitialSalepersonIDValue() → string
```

**Parámetros de búsqueda:**
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

### 2. **Actividades** (Activities)
```
Entidad: XPlorActivity
├── ActivityID
├── Name (ej: "PADEL", "Fitness")
├── ClubID
├── ColorHex (color de la actividad)
├── Durations[] (ISO 8601: PT60M, PT90M)
├── IsBookable, IsViewable
├── ShowcaseActivities[] (actividades de vitrina)
├── ActivityGroups[] (grupos de actividades)
├── ArchivedAt/By (fecha y autor de archivado)
├── CreatedAt/By
└── TemplateToken

Métodos útiles:
- ActivityID() → string
- ShowcaseIDs() → []int (IDs de actividades en vitrina)
- IsActive() → bool (no archivada)
- IsPadel() → bool (verifica si es pádel)
- DurationMinutes() → int (extrae minutos de PT format)
```

**Parámetros de búsqueda:**
```
XPlorActivitiesParams:
- ClubID / ClubIDs[]
- Name
- Archived (bool)
```

---

### 3. **Clases** (Class Events)
```
Entidad: XPlorClass
├── ClassEventID
├── Club, Studio, Activity, Coach
├── StartedAt, EndedAt (LocalTime)
├── AttendingLimit, QueueLimit
├── BookedAttendees[] (asistentes confirmados)
├── QueuedAttendees[] (asistentes en lista de espera)
├── AttendeeRemaining (plazas disponibles)
├── QueueRemaining (plazas de espera)
├── Summary, Description
├── CoachAvailable (bool)
├── DisabledItems[] (ítems deshabilitados)
├── Recurrence (para clases recurrentes)
├── ClassLayout, ClassLayoutConfiguration
├── Processing (bool)
└── CreatedAt/UpdatedAt/DeletedAt/ArchivedAt

Métodos útiles:
- ClassEventID() → string
- ClubID() → string
- StudioID() → string
- ActivityID() → string
- CoachID() → string
- RecurrenceID() → string
- GetStartedAt() → time.Time
- GetEndedAt() → time.Time
- HasAvailableSpots() → bool
- HasQueueSpots() → bool
- IsActive() → bool (no borrada ni archivada)
- IsDeleted() → bool
- GetAllContactIDs() → []string (IDs de todos los asistentes)
```

**Parámetros de búsqueda:**
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

### 4. **Estudios** (Studios)
```
Entidad: XPlorStudio
├── StudioID
├── Name
├── Club (club al que pertenece)
├── ZoneID (zona dentro del club)
├── Capacity (capacidad total)
├── Overbooking (porcentaje de sobreventa)
├── StreetAddress, PostalCode
├── AddressLocality, AddressCountry
├── Tags
├── CreatedAt/By
└── ArchivedAt/By

Métodos útiles:
- StudioID() → string
- ClubID() → string
- ZoneID() → string
- Address() → string (dirección completa)
```

---

### 5. **Entrenadores** (Coaches)
```
Entidad: XPloreCoach
├── CoachID
├── GivenName, FamilyName
├── AlternateName
├── Email, Mobile
├── Activities[] (actividades que enseña)
├── CreatedAt/By
└── ArchivedAt/By

Métodos útiles:
- CoachID() → string
- ActivityIDs() → []string (IDs de las actividades)
```

---

### 6. **Clubes** (Clubs)
```
Entidad: XPlorClub
├── ClubID, ClubNumberID (id numérico)
├── Code (3-5 caracteres)
├── Name
├── Number
├── Email, Phone
├── StreetAddress, PostalCode
├── AddressLocality, AddressCountry
├── AddressCountryIso
├── OpeningDate
├── Description
├── ClubTags[] (etiquetas del club)
├── PublicMetadata (metadatos públicos)
├── SaleTerms[] (términos de venta)
├── Locale
├── CreatedAt/By
└── DeletedAt

Campos relacionados:
- TaxRates (impuestos)
- ResaboxNotification (notificaciones)
```

---

### 7. **Suscripciones** (Subscriptions)
```
Entidad: XPlorSubscription
├── SubscriptionID
├── Contact (contacto/cliente)
├── Club (club de la suscripción)
├── StartedAt, EndedAt
├── IsActive (bool)
├── CreatedAt/UpdatedAt
└── [Otros campos específicos del plan]

Métodos útiles:
- SubscriptionID() → string
- ContactID() → string
- ClubID() → string
- IsActive() → bool
- IsExpired() → bool
```

**Parámetros de búsqueda:**
```
XPlorSubscriptionsParams:
- ContactID / ContactIDs[]
- ClubID
- Active (bool)
- StartedAt / EndedAt (rango de fechas)
```

---

### 8. **Usuarios del Sistema** (Users)
```
Entidad: XPlorUser
├── UserID
├── Email
├── GivenName, FamilyName
├── Mobile, PictureLink
├── Code
├── ClubIds[] (clubes a los que tiene acceso)
├── NetworkNodeIds[] (nodos de red)
├── Roles[] (roles del usuario)
├── Active (bool)
├── Locale
├── CreatedAt/DeletedAt/ArchivedAt
└── [Otros campos]

Métodos útiles:
- UserID() → string
- ClubIDs() → []string
- NetworkNodeIDs() → []string
- PropertiesNetworkNodeIDs() → []string
- IsActive() → bool
- IsDeleted() → bool
- IsArchived() → bool
- IsInactive() → bool
- FullName() → string
- HasRole(role string) → bool
- GetCreatedAt() → time.Time
- GetDeletedAt() → *time.Time
```

---

### 9. **Nodos de Red** (Network Nodes)
```
Entidad: XPlorNetworkNode
├── NodeID
├── Name
├── Type
├── RelatedClubs[] (clubes relacionados)
└── [Otros campos]
```

---

### 10. **Eventos** (Events)
```
Entidad: XPlorEvent
├── EventID
├── Name
├── StartAt, EndAt
├── Location
└── [Otros campos]
```

---

### 11. **Familias** (Families)
```
Entidad: XPlorFamily
├── FamilyID
├── Members[] (miembros de la familia)
├── CreatedAt/UpdatedAt
└── [Otros campos]

Parámetros de búsqueda:
XPlorFamiliesParams:
- FamilyID / FamilyIDs[]
- ContactID
```

---

### 12. **Otras Entidades**

- **Recurrencias**: Patrones de repetición para clases (diaria, semanal, etc.)
- **Tipos de Clase**: Tipos de clases disponibles
- **Contadores**: Líneas de contador/efectivo
- **Etiquetas de Contacto**: Tags para categorizar contactos
- **Artículos**: Contenido/artículos del sistema
- **Zonas**: Áreas dentro de un club

---

## 📡 Funciones Públicas del Provider

### Gestión de Contactos
```go
Contacts(nodeId string, params *XPlorContactsParams, 
         pagination *XPlorPagination) → (*XPlorContacts, error)
Contact(nodeId string, contactId string) → (*XPlorContact, error)
```

### Gestión de Actividades
```go
Activities(nodeId string, queryParams *XPlorActivitiesParams,
          pagination *XPlorPagination) → (*XPlorActivities, error)
Activity(nodeId string, activityId string) → (*XPlorActivity, error)
```

### Gestión de Clases
```go
Classes(nodeId string, params *XPlorClassesParams,
       pagination *XPlorPagination) → (*XPlorClasses, error)
Class(nodeId string, classId string) → (*XPlorClass, error)
```

### Gestión de Estudios
```go
Studios(nodeId string, pagination *XPlorPagination) → (*XPlorStudios, error)
Studio(nodeId string, studioId string) → (*XPlorStudio, error)
```

### Gestión de Entrenadores
```go
Coaches(nodeId string, pagination *XPlorPagination) → (*XPloreCoaches, error)
Coach(nodeId string, coachId string) → (*XPloreCoach, error)
```

### Gestión de Clubes
```go
Clubs(nodeId string) → (*XPloreClubs, error)
Club(nodeId string, clubId string) → (*XPlorClub, error)
```

### Gestión de Suscripciones
```go
Subscriptions(nodeId string, params *XPlorSubscriptionsParams,
             pagination *XPlorPagination) → (*XPlorSubscriptions, error)
Subscription(nodeId string, subscriptionId string) → (*XPlorSubscription, error)
```

### Gestión de Usuarios
```go
Users(pagination *XPlorPagination) → (*XPlorUsers, error)
User(userId string) → (*XPlorUser, error)
```

### Gestión de Nodos de Red
```go
NetworkNodes(pagination *XPlorPagination) → (*XPlorNetworkNodes, error)
NetworkNode(nodeId string) → (*XPlorNetworkNode, error)
```

### Gestión de Asistentes
```go
Attendees(nodeId string, classId *string,
         pagination *XPlorPagination) → (*XPlorAttendees, error)
```

### Gestión de Eventos
```go
Events(nodeId string, pagination *XPlorPagination,
      timeGap *XPlorTimeGap) → (*XPlorEvents, error)
```

### Gestión de Familias
```go
Families(nodeId string, params *XPlorFamiliesParams,
        pagination *XPlorPagination) → (*XPlorFamilies, error)
Family(nodeId string, familyId string) → (*XPlorFamily, error)
```

### Gestión de Recurrencias
```go
Recurrences(nodeId string, params *XPlorRecurrencesParams,
           pagination *XPlorPagination) → (*XPlorRecurrences, error)
Recurrence(nodeId string, recurrenceId string) → (*XPlorRecurrence, error)
```

### Gestión de Etiquetas de Contacto
```go
ContactTags(nodeId string, params *XPlorContactTagsParams,
           pagination *XPlorPagination) → (*XPlorContactTags, error)
ContactTag(nodeId string, contactTagId string) → (*XPlorContactTag, error)
```

### Gestión de Tipos de Clase
```go
ClassType(nodeId string, classTypeId string) → (*XPlorClassType, error)
```

### Gestión de Contadores
```go
CounterLines(nodeId string, pagination *XPlorPagination) → (*XPlorCounterLines, error)
CounterLine(nodeId string, counterLineId string) → (*XPlorCounterLine, error)
```

### Gestión de Artículos
```go
Articles(nodeId string, pagination *XPlorPagination) → (*XPlorArticles, error)
Article(nodeId string, articleId string) → (*XPlorArticle, error)
```

### Gestión de Zonas
```go
Zones(nodeId string, params *XPlorZonesParams,
     pagination *XPlorPagination) → (*XPlorZones, error)
Zone(nodeId string, zoneId string) → (*XPlorZone, error)
```

### Cierre
```go
Close() → void (libera recursos del provider)
```

---

## 📍 Parámetros Comunes

### Paginación
```go
type XPlorPagination struct {
    Page         int // Número de página (1-indexed)
    ItemsPerPage int // Elementos por página
}
```

### Rango Temporal
```go
type XPlorTimeGap struct {
    StartAt *time.Time // Fecha/hora de inicio
    EndAt   *time.Time // Fecha/hora de fin
}
```

### Construcción de Parámetros de Query
```go
// Para paginación
BuildPaginationQueryParams(pagination) → url.Values

// Para paginación + rango temporal
BuildPaginationAndTimeGapParams(pagination, timeGap) → url.Values
```

---

## 🔄 Patrón de Uso General

```go
// 1. Inicializar el provider
provider := xplorcore.Init(config)

// 2. Llamar a las funciones públicas con nodeId
contacts, err := provider.Contacts(nodeId, params, pagination)
if err != nil {
    // Manejar error
}

// 3. Trabajar con los resultados
for _, contact := range contacts.Contacts {
    contactID, _ := contact.ContactID()
    fmt.Println(contact.GivenName, contact.FamilyName, contactID)
}

// 4. Cerrar cuando termine
provider.Close()
```

---

## 🔒 Características de Seguridad

1. **Autenticación OAuth2**: Credenciales seguras
2. **Token Caching**: Reutiliza tokens válidos
3. **Thread-Safety**: Sincronización automática con mutex
4. **Pool de Ejecutores**: Reutiliza conexiones HTTP
5. **Headers Dinámicos**: Inyecta contexto de usuario (Node ID, Club ID)

---

## 📝 Notas de Implementación

- Las funciones retornan `(*Entity, error)` o `(*EntityCollection, error)`
- El parámetro `nodeId` es requerido en la mayoría de funciones (identifica la ubicación)
- Algunos métodos no necesitan `nodeId` (ej: `Users`, `NetworkNodes`)
- Los IDs se extraen automáticamente de las rutas IRI (ej: "/enjoy/clubs/1249" → "1249")
- Todas las fechas usan el formato `LocalTime` o `LocalDate` con parsing automático
- El SDK maneja sincronización automática entre llamadas concurrentes

---

## 🎭 Ejemplo Completo

```go
package main

import (
    "github.com/angelbarreiros/XPlorGo/xplorcore"
    "github.com/angelbarreiros/XPlorGo/xplorentities"
)

func main() {
    // Inicializar
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

    // Obtener actividades con paginación
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

    // Procesar actividades
    for _, activity := range activities.Activities {
        println(activity.Name)
    }

    // Cerrar
    provider.Close()
}
```

---


