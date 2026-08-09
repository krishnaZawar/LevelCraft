# Input Handling and Modelling Design

## 1. Overview

This document defines how user input is modelled, received, normalized, stored, and consumed by the game Builder.

The design intentionally focuses on the **input pipeline** and its integration with the existing Builder/interpreter architecture. 

The central principles are:

- Input from the UI is **transition-based**, not frame-based.
- The backend is the **source of truth** for the canonical input registry.
- Input identifiers have a stable **symbolic name** and an internal numeric code.
- The UI obtains the backend's input mapping through an explicit API.
- A `CommandRequest` represents an input transition such as press or release.
- The backend maintains persistent input state independently of the UI's frame rate.
- Components own their input bindings because they own the actions associated with those bindings.
- Interactable components are identified explicitly so the interpreter can efficiently filter potential input consumers.
- Every matching component/object may independently emit an event for the same input.
- Events emitted by input handling operate on the object that generated them; interactions with other objects are resolved during event processing.

---

## 2. Goals

### 2.1 Primary goals

The input system should:

1. Avoid sending the complete input state every UI frame.
2. Represent key/button transitions exactly once per transition.
3. Maintain authoritative input state in the backend.
4. Provide a stable input representation between UI and backend.
5. Allow the backend to remain independent of the UI implementation.
6. Allow input bindings to be stored directly on component instances.
7. Allow users to change bindings without changing interpreter logic.
8. Allow multiple objects/components to respond to the same input.
9. Provide an extensible model for additional input devices in the future.
10. Keep the hot path efficient enough for the expected number of game objects.

### 2.2 Non-goals

The following are outside the scope of this document:

- Detailed event-processing architecture.
- Physics and AI simulation.
- Game-specific behavior.
- Component execution semantics.
- Cross-object event processing.
- Runtime registration of new input types.
- Detailed continuous-input semantics.

---

# 3. High-Level Architecture

The input pipeline is:

```text
┌─────────────────────┐
│         UI          │
│                     │
│ Native input event  │
└──────────┬──────────┘
           │
           │ normalize using
           │ backend-provided schema
           ▼
┌─────────────────────┐
│  CommandRequest     │
│                     │
│ Input + transition  │
│ + metadata          │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ CommandRequestQueue │
└──────────┬──────────┘
           │
           │ game-loop tick
           ▼
┌─────────────────────┐
│    Input Manager    │
│                     │
│ Current state       │
│ Pressed transitions │
│ Released transitions│
└──────────┬──────────┘
           │
           ├───────────────────────┐
           │                       │
           ▼                       ▼
┌─────────────────────┐   ┌─────────────────────┐
│    Interpreter      │   │     Simulation      │
│                     │   │                     │
│ Process transitions │   │ Consume input state │
└──────────┬──────────┘   └─────────────────────┘
           │
           ▼
┌─────────────────────┐
│       Events        │
└─────────────────────┘
```

The important separation is:

> **The UI reports input transitions. The backend owns input state.**

The UI therefore does not need to repeatedly report that a key remains held.

---

# 4. Input Representation

## 4.1 Canonical input identifier

The engine should define a canonical `InputCode`.

The numeric value is an internal engine identifier. It may use ASCII-compatible values for printable keys, but ASCII is **not** the conceptual model of the input system.

For example:

```text
KEY_A       → 65
KEY_W       → 87
KEY_SPACE   → 32
```

Special inputs can occupy an engine-defined range.

The exact numeric values are owned by the backend and are not required to have semantic meaning.

### Why numeric codes?

A compact numeric identifier is useful for runtime operations:

```text
inputState[87]
```

can directly access the state of `KEY_W`.

This avoids string lookups in the input hot path.

### Why not treat the number as the public semantic identifier?

Numeric values may change as the engine evolves. A stable symbolic name is therefore retained as the externally meaningful identifier.

The important invariant is:

```text
KEY_W
```

continues to identify the same logical engine input even if its internal numeric representation changes.

---

# 5. Backend Input Registry

The backend owns a static input registry.

Conceptually:

```text
Input Registry

Symbolic Name        Numeric Code
----------------------------------
KEY_A                65
KEY_B                66
...
KEY_W                87
KEY_SPACE            32
KEY_ESCAPE           ...
MOUSE_LEFT           ...
MOUSE_RIGHT          ...
```

The registry is:

- backend-owned;
- authoritative;
- immutable at runtime;
- extensible by engine development;
- independent of component instances.

Runtime registration is intentionally not supported in the current design.

---

# 6. Input Schema API

The UI should not maintain a duplicated input-definition file.

Instead, the backend exposes the registry through an API.

Example:

```http
GET /api/v1/input/schema
```

Example response:

```json
{
  "version": 1,
  "inputs": [
    {
      "name": "KEY_W",
      "code": 87,
      "uiCode": "KeyW",
      "label": "W"
    },
    {
      "name": "KEY_SPACE",
      "code": 32,
      "uiCode": "Space",
      "label": "Space"
    },
    {
      "name": "KEY_ESCAPE",
      "code": 256,
      "uiCode": "Escape",
      "label": "Escape"
    },
    {
      "name": "MOUSE_LEFT",
      "code": 400,
      "uiCode": "MouseLeft",
      "label": "Left Mouse Button"
    }
  ]
}
```

The exact schema can evolve, but the conceptual responsibilities are:

| Field | Purpose |
|---|---|
| `name` | Stable symbolic engine identifier |
| `code` | Backend numeric runtime identifier |
| `uiCode` | Mapping from the UI's native input representation |
| `label` | Human-readable representation |
| `version` | Schema/version compatibility |

## 6.1 Why the backend is the source of truth

Keeping the registry in the backend avoids maintaining synchronized definitions across repositories.

This is particularly important if the UI and backend are later split into separate repositories.

The dependency becomes:

```text
UI
 │
 │ asks
 ▼
Backend Input Schema
```

rather than:

```text
Shared File
   ├── Backend
   └── UI
```

This also allows the backend to evolve its internal numeric representation without requiring the UI to understand how those numbers are assigned.

## 6.2 UI initialization

The UI should obtain the input schema during initialization:

```text
UI starts
   ↓
GET /input/schema
   ↓
Store schema in UI runtime state
   ↓
Use schema for input translation
   ↓
Send CommandRequests
```

The schema can be cached for the lifetime of the Builder session.

---

# 7. UI-to-Backend Input Translation

The UI has its own native input representation.

For example, a browser may provide:

```text
KeyboardEvent.code = "KeyW"
```

The UI translates this through the schema obtained from the backend:

```text
"KeyW"
    ↓
KEY_W
    ↓
87
```

The UI then creates a `CommandRequest`.

The backend therefore does not need to depend on browser-specific input representations.

The boundary is:

```text
UI-native input
       ↓
Input schema
       ↓
Engine InputCode
       ↓
CommandRequest
```

This also keeps the engine independent from the particular UI technology.

---

# 8. CommandRequest

A `CommandRequest` represents a **single input transition** received from the UI.

It should not represent the entire input state of a frame.

Where `CommandType` can contain transition types such as:

```text
KEY_PRESSED
KEY_RELEASED
BUTTON_PRESSED
BUTTON_RELEASED
```

The exact device terminology is intentionally abstracted so that the model can be extended later.

---

# 9. Transition-Based Input

The UI should send exactly one request for each relevant transition.

For example, when the user holds `W`:

```text
W pressed
    ↓
KEY_PRESSED(W)
```

No additional request is sent while the key remains held.

When the user releases `W`:

```text
W released
    ↓
KEY_RELEASED(W)
```

The communication is therefore:

```text
UI                              Backend

W DOWN ───────────────────────→ KEY_PRESSED(W)

          [held for 2 seconds]

W UP   ───────────────────────→ KEY_RELEASED(W)
```

Not:

```text
Frame 1 → W
Frame 2 → W
Frame 3 → W
...
Frame N → W
```

and not:

```text
Frame 1 → [W, A, Space, ...]
Frame 2 → [W, A, Space, ...]
...
```

This keeps UI communication independent of simulation frequency.

---

# 10. CommandRequest Queue

All incoming `CommandRequest`s are placed into the existing command queue.

Example:

```text
UI
 │
 ├── KEY_PRESSED(W)
 ├── KEY_PRESSED(MOUSE_LEFT)
 ├── KEY_RELEASED(MOUSE_LEFT)
 └── KEY_RELEASED(W)
        │
        ▼
CommandRequestQueue
```

The queue preserves request ordering.

---

# 11. Input Manager

The Input Manager converts queued input transitions into authoritative backend input state.

It maintains three conceptual states:

```text
Current / Down
Pressed
Released
```

For example:

```text
Current:
    W = DOWN

Pressed this tick:
    W

Released this tick:
    none
```

After the user releases W:

```text
Current:
    W = UP

Pressed this tick:
    none

Released this tick:
    W
```

A possible representation is:

```go
type InputState struct {
    Down     []bool
    Pressed  []bool
    Released []bool
}
```

`Down` persists across ticks.

`Pressed` and `Released` represent transitions observed during the current processing interval and are cleared after the relevant tick.

---

# 12. Example: Holding W

### Step 1: UI

The user presses W.

```text
KEY_PRESSED(W)
```

is sent once.

### Step 2: Queue

```text
CommandRequestQueue
    ↓
KEY_PRESSED(W)
```

### Step 3: Input Manager

```text
Down[KEY_W] = true
Pressed[KEY_W] = true
```

### Step 4: Subsequent ticks

No further UI request is necessary.

The backend still has:

```text
Down[KEY_W] = true
```

### Step 5: UI release

The user releases W:

```text
KEY_RELEASED(W)
```

### Step 6: Input Manager

```text
Down[KEY_W] = false
Released[KEY_W] = true
```

This allows the rest of the engine to distinguish:

```text
Was W pressed?
Is W currently held?
Was W released?
```

without requiring repeated network messages.

---

# 13. Ordering of Press and Release

Input transitions must be processed in the order in which they were received.

Consider:

```text
101 → KEY_PRESSED(W)
102 → KEY_RELEASED(W)
103 → KEY_PRESSED(SPACE)
```

The Input Manager applies them in sequence.

This preserves the actual input history even when multiple transitions occur before a game-loop tick.

This is particularly important for short-lived inputs.

For example:

```text
W pressed
W released
```

may both occur between two simulation ticks.

A simple current-state-only representation could lose the fact that the transition happened.

The transient `Pressed` and `Released` state preserves that information.

---

# 14. Focus and Connection Loss

A UI may lose focus before sending a release event.

For example:

```text
W pressed
    ↓
KEY_PRESSED(W)

Browser/tab loses focus
    ↓
KEY_RELEASED(W) never arrives
```

The backend could incorrectly retain:

```text
W = DOWN
```

To avoid this, the input protocol should support an input-reset/lifecycle signal.

For example:

```text
INPUT_FOCUS_LOST
```

or:

```text
INPUT_RESET
```

When received, the Input Manager clears active inputs:

```text
W = UP
A = UP
MouseLeft = UP
...
```

The same reset behavior should be triggered when the input connection terminates unexpectedly.

---

# 15. Input Bindings Belong to Component Instances

Input mappings that define component behavior should be stored on the component instance itself.

Example:

```text
Player
 └── WeaponComponent
       └── ShootInput = MOUSE_LEFT
```

The user can change it:

```text
ShootInput = MOUSE_RIGHT
```

without any change to the interpreter.

The component owns the action, therefore it owns the input binding associated with that action.

The binding should use the stable symbolic input identity rather than depending on an incidental numeric value.

Conceptually:

```json
{
  "component": "WeaponComponent",
  "bindings": {
    "shoot": "MOUSE_RIGHT"
  }
}
```

At runtime, the engine resolves the symbolic identifier to the current canonical `InputCode`.

This allows internal numeric assignments to evolve without invalidating the semantic component configuration.

---

# 16. Why Input Bindings Are Not Global

A global mapping such as:

```text
MOUSE_LEFT → SHOOT
```

would introduce game-specific assumptions into the engine.

Different objects can legitimately use the same input differently.

For example:

```text
Player A
    MOUSE_LEFT → Shoot

Player B
    MOUSE_LEFT → Interact

UI Button
    MOUSE_LEFT → Activate
```

The interpreter should therefore not know what `MOUSE_LEFT` means globally.

It should only determine:

> Which interactable component instances have a binding matching this input transition?

The component then emits its configured event.

---

# 17. Interactable Components

Components that can respond to external input should expose an `Interactable` capability/marker.

The interpreter can therefore perform an initial filter:

```text
All game objects
      ↓
Interactable components
      ↓
Check input bindings
      ↓
Matching components
      ↓
Emit events
```

This avoids inspecting irrelevant components.

The `Interactable` property is intentionally generic.

It does not encode what the component does.

---

# 18. Generic Interpreter Input Flow

When an input transition reaches the interpreter:

```text
CommandRequest
      ↓
Input transition
      ↓
Iterate objects
      ↓
Find interactable components
      ↓
Check component input bindings
      ↓
Binding matches?
      │
      ├── No → continue
      │
      └── Yes
            ↓
        Component emits
        its configured event
```

The interpreter does not contain mappings such as:

```text
W → Move
LeftClick → Shoot
Space → Jump
```

Those mappings belong to component instances.

This is what allows the same interpreter to support different game types.

---

# 19. Multiple Components May Respond to the Same Input

A single input transition may match multiple component instances.

For example:

```text
MOUSE_LEFT pressed
```

could match:

```text
Player
 └── WeaponComponent
       └── Shoot = MOUSE_LEFT

Door
 └── InteractionComponent
       └── Activate = MOUSE_LEFT

Button
 └── ButtonComponent
       └── Activate = MOUSE_LEFT
```

Each matching component independently emits its event.

The interpreter does not select one "primary" consumer.

Conceptually:

```text
MOUSE_LEFT
    │
    ├── WeaponComponent → ShootEvent
    ├── InteractionComponent → ActivateEvent
    └── ButtonComponent → ActivateEvent
```

Whether those events are ultimately processed, ignored, or cause state changes is the responsibility of the existing event-processing architecture.

---

# 20. Event Ownership

An event emitted from input handling operates on the object whose component generated the event.

Example:

```text
Player
 └── WeaponComponent
       └── MOUSE_LEFT → ShootEvent
```

The resulting event is associated with:

```text
Player
```

The input system does not directly resolve another target object.

If processing the event requires another object, that resolution occurs during event processing.

This keeps input interpretation independent from game-specific object relationships.

---

# 21. Input Schema and Component Configuration

The complete relationship is:

```text
Backend Input Registry
        │
        │ GET /input/schema
        ▼
       UI
        │
        │ user chooses input
        ▼
Component Configuration
        │
        │ symbolic input name
        ▼
    Saved Scene
        │
        ▼
Component Instance
        │
        │ input binding
        ▼
    Interpreter
```

For example:

```text
Backend Registry:

KEY_W       → 87
MOUSE_LEFT  → 400
MOUSE_RIGHT → 401
```

User configures:

```text
WeaponComponent
    Shoot = MOUSE_RIGHT
```

The interpreter later receives:

```text
MOUSE_RIGHT pressed
```

and performs a generic binding lookup.

It does not need to know anything about shooting.

---

# 22. Runtime Representation

The symbolic name is useful for:

- scene persistence;
- configuration;
- debugging;
- API responses;
- UI display.

The numeric `InputCode` is useful for runtime processing.

Therefore the system can use:

```text
Persistent/configuration layer:
    "MOUSE_RIGHT"

Protocol/runtime layer:
    InputCode(401)
```

The backend resolves the symbolic identifier when loading/configuring the component.

This gives the system both stability and efficient runtime access.

---

# 23. Future Optimization

The initial interpreter can use a straightforward scan:

```text
for each object:
    for each interactable component:
        if binding matches input:
            emit event
```

This is acceptable while game scenes contain relatively few objects.

If the number of interactable objects grows significantly, the engine can introduce an index:

```text
InputCode
    ↓
Interactable component instances
```

For example:

```text
MOUSE_LEFT
    ├── Player.WeaponComponent
    ├── Door.InteractionComponent
    └── Button.ActivationComponent

KEY_SPACE
    ├── Player.JumpComponent
    └── Character2.JumpComponent
```

The input pipeline does not need to change.

Only the lookup strategy changes:

```text
Current:
Input → scan interactables

Future:
Input → indexed interactables
```

This is why component bindings and the interpreter should remain decoupled from the indexing strategy.

---

# 24. Design Invariants

The following invariants should be maintained.

### Input transport

- One relevant press transition produces one `CommandRequest`.
- One relevant release transition produces one `CommandRequest`.
- Held inputs do not generate repeated network requests.
- Input ordering is preserved.

### Input authority

- The backend owns authoritative input state.
- The UI does not own the simulation's input state.
- The backend input registry is the source of truth.

### Input identity

- Symbolic input names represent stable semantic identities.
- Numeric codes are runtime identifiers.
- Numeric assignments may change without changing the semantic name.

### Component ownership

- A component owns the input binding for its action.
- The interpreter does not contain game-specific input mappings.
- Multiple component instances may respond to the same input.

### Event generation

- Matching interactable components independently emit events.
- Input handling does not directly perform game-specific behavior.
- Events initially operate on the object that generated them.

### Extensibility

- The input model is device agnostic.
- New input categories can be added without changing the fundamental `CommandRequest → InputManager → Interpreter` pipeline.
- Runtime input registration is not supported in the current version.

---

# 25. Summary

The proposed input architecture can be reduced to four responsibilities:

```text
UI
    ↓
Report transitions

Input Manager
    ↓
Maintain authoritative state

Interpreter
    ↓
Match transitions against component-owned bindings

Event System
    ↓
Execute the resulting behavior
```

The most important design decision is that **input transport and input state are separate concepts**.

The UI says:

```text
"MOUSE_LEFT was pressed."
```

The backend maintains:

```text
"MOUSE_LEFT is currently down."
```

The interpreter asks:

```text
"Which interactable component instances are configured to respond to MOUSE_LEFT?"
```

Those components emit their respective events.

This provides a device-agnostic, game-agnostic input pipeline while keeping input configuration local to the components that own the behavior.
