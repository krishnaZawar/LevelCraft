# Builder Architecture Design Document

## 1. Overview

This document defines the architecture and runtime lifecycle of the
**Builder** component of a generic no-code game engine.

The Builder is responsible for:

-   Receiving user input from the UI.
-   Converting external input into generic command requests.
-   Queuing commands for deterministic game-loop processing.
-   Interpreting commands against the current game state.
-   Discovering components capable of reacting to the command.
-   Allowing those components to generate task-specific events.
-   Processing events and producing further state-changing events.
-   Running continuous simulations such as AI, physics, movement, and
    collision detection.
-   Applying resulting state changes.
-   Exposing the updated game state to the UI for rendering.

The Builder must remain **game-type agnostic**.

Instead, the Builder understands generic concepts:

-   Commands.
-   Entities.
-   Components.
-   Interactable components.
-   Events.
-   Systems/simulations.
-   State transitions.
-   Entity references.
-   Context.

The behavior of a particular game emerges from the combination of
components and their metadata.

------------------------------------------------------------------------

# 2. Design Goals

## 2.1 Game-Type Agnostic Core

The Builder must not contain hardcoded knowledge such as:

``` text
LEFT CLICK means shoot.
E means interact.
Enemy means hostile.
Weapon means projectile.
```

Those are game-specific semantics.

The Builder should instead understand:

``` text
An input command occurred.
A component may be capable of reacting to it.
The component can generate one or more events.
Events can cause further state changes.
Simulation systems run continuously.
```

------------------------------------------------------------------------

## 2.2 Component-Oriented Behavior

A component ideally represents one primary capability or task.

For example:

``` text
WeaponComponent
    → handles weapon-related action

MovementComponent
    → handles movement

HealthComponent
    → handles health state

CollisionComponent
    → participates in collision detection

PickupComponent
    → handles pickup behavior
```

More complex behavior should be created by composing components.

For example:

``` text
Player
├── TransformComponent
├── WeaponComponent
├── MovementComponent
├── CollisionComponent
└── HealthComponent
```

The Builder should not need to know that this combination represents a
"player".

The behavior emerges from the components.

------------------------------------------------------------------------

## 2.3 Input Remapping Must Not Break Behavior

Physical input and game behavior must be separate.

For example:

``` text
LEFT CLICK → Weapon action
```

may later become:

``` text
RIGHT CLICK → Weapon action
```

without changing the WeaponComponent.

The component should not be permanently coupled to a physical input such
as:

``` text
LEFT_CLICK
```

Instead, the input mapping is configuration.

The component exposes its interaction/task configuration independently.

The exact representation may vary, but the architectural requirement is:

``` text
Physical Input
    ↓
Input Mapping / Matching
    ↓
Interactable Component
    ↓
Task Event
```

The component remains valid if the physical input configuration changes.

------------------------------------------------------------------------

## 2.4 Deterministic Runtime

The Builder should provide a deterministic lifecycle:

``` text
Input Ingestion
    ↓
Command Processing
    ↓
Event Processing
    ↓
Simulation
    ↓
Collision/Detection
    ↓
Deferred Events
    ↓
State Update
    ↓
Render Snapshot
```

The exact event timing may vary, but the lifecycle must be well-defined.

------------------------------------------------------------------------

## 2.5 Extensibility

The initial implementation may use straightforward scanning:

``` text
For every entity:
    For every component:
        If component is interactable:
            Check whether it matches the command
```

This is acceptable for the initial scale of a typical 2D game.

The architecture should later allow indexing and graph-based dispatch
without changing the external component model.

Possible future optimizations:

``` text
All Entities
    ↓
Interactable Components
    ↓
Components Registered for Input Type
    ↓
Components Matching Input
    ↓
Context-Compatible Components
```

------------------------------------------------------------------------

# 3. High-Level Architecture

``` text
┌─────────────────────┐
│       UI            │
│                     │
│ Keyboard / Mouse    │
│ Gamepad / Editor    │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  CommandRequest     │
│                     │
│ Generic external    │
│ input representation│
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ CommandRequestQueue │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│     Game Loop       │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│    Interpreter      │
│                     │
│ Finds components    │
│ capable of reacting │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│      Events         │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│   Event Processor   │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  State Operations   │
│  / More Events      │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│    Simulation       │
│                     │
│ AI / Physics /      │
│ Movement / Timers   │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ Detection Systems   │
│                     │
│ Collision / Trigger │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ Deferred Events     │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ Updated Game State  │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  Render Snapshot    │
│        → UI         │
└─────────────────────┘
```

------------------------------------------------------------------------

# 4. Input Mapping

Physical input must be decoupled from component behavior.

A possible conceptual mapping is:

``` text
Physical Input
    ↓
Matching Configuration
    ↓
Interactable Component
    ↓
Task Event
```

Example:

``` text
LEFT_CLICK
    ↓
WeaponComponent
    ↓
ShootEvent
```

If the user changes the configuration:

``` text
RIGHT_CLICK
    ↓
WeaponComponent
    ↓
ShootEvent
```

the WeaponComponent does not need to change.

The exact implementation can be direct matching rather than an explicit
semantic action layer.

The component can internally determine whether the incoming request
matches its configured trigger.

The Interpreter remains unaware of the task.

------------------------------------------------------------------------

# 5. The Interpreter

The Interpreter is the central mechanism that converts external commands
into game events.

Its responsibilities are:

1.  Receive a `CommandRequest`.
2.  Iterate through candidate entities/components.
3.  Identify components that can interact.
4.  Determine whether the component accepts the command.
5.  Resolve the interaction context.
6.  Ask the component to produce its task event(s).
7.  Add those events to the event queue.

The Interpreter must not contain logic such as:

``` text
if input == LEFT_CLICK:
    shoot()
```

or:

``` text
if component is Weapon:
    create ShootEvent
```

Instead:

``` text
if component implements Interactable:
    if component matches request:
        events = component.Interact(...)
```

------------------------------------------------------------------------

## 5.1 Initial Interpreter Algorithm

The initial implementation can be:

``` text
function Interpret(command):
    events = []
    for each entity in world:
        for each component in entity:
            if component is not Interactable:
                continue
            context = resolve context(command, entity)
            if component does not match command:
                continue
            generatedEvents = component.Interact(command, context)
            events.append(generatedEvents)
    return events
```

------------------------------------------------------------------------

# 6. Why Scanning Is Acceptable Initially

The initial complexity is approximately:

``` text
O(number of entities × components per entity)
```

Example:

``` text
1,000 entities × 5 components per entity = 5,000 component checks
```

This is generally acceptable for an initial implementation, especially
because:

-   The scan happens only when a command is processed.
-   Most components are not interactable.
-   The check can be a simple interface/type check.
-   The game loop can process commands at a fixed point in the frame.

The architecture should not be prematurely optimized for millions of
entities.

The first priority is correctness and extensibility.

------------------------------------------------------------------------

# 7. Interaction Context

A command by itself is usually insufficient to determine the behavior.

For example:

``` text
LEFT_CLICK at (800, 450)
```

may require the Builder to determine:

-   Which entity owns the interaction?
-   What entity is under the cursor?
-   What is the world-space position?
-   What is the screen-space position?
-   What is the current camera?
-   What is the current selection?
-   What is the current game state?
-   Is the target valid?

Therefore the Interpreter should construct an `InteractionContext`.

The context should be treated as a controlled view of the world.

The component should not receive unrestricted access to every mutable
object unless there is a specific reason.

------------------------------------------------------------------------

# 8. Components Generate Task Events

The component owns the task semantics.

For example:

``` text
WeaponComponent
    → ShootEvent

DoorInteractionComponent
    → OpenDoorEvent

PickupComponent
    → PickupEvent
```

The Interpreter does not know these event types.

It simply asks the component:

``` text
Interact(...)
```

and receives:

``` text
[]Event
```

------------------------------------------------------------------------

# 9. Example: Top-Down Shooter

Consider a top-down shooter.

## 9.1 Scene

``` text
Player
├── TransformComponent
├── WeaponComponent
└── HealthComponent

Enemy
├── TransformComponent
├── AIComponent
├── HealthComponent
└── CollisionComponent

Bullet
├── TransformComponent
├── MovementComponent
├── CollisionComponent
└── DamageComponent
```

The Builder does not need to know that the first entity is a player.

The component combination produces the behavior.

## 9.2 User Input

The user clicks:

``` text
LEFT MOUSE BUTTON
at (800, 450)
```

The UI sends:

``` json
{
  "type": "POINTER_BUTTON_PRESSED",
  "button": "LEFT",
  "position": {
    "x": 800,
    "y": 450
  }
}
```

The command enters:

``` text
CommandQueue
```

## 9.3 Game Loop Reads the Command

The game loop dequeues:

``` text
PointerButtonPressed
```

and passes it to:

``` text
Interpreter
```

## 9.4 Interpreter Scans the World

The Interpreter examines:

``` text
Player
    TransformComponent
    WeaponComponent
    HealthComponent
```

It identifies:

``` text
WeaponComponent
    implements Interactable
```

It checks whether the command matches the component configuration.

The component generates:

``` text
ShootEvent
```

The event is added to the event queue.

# 10. Event Processing

The event processor reads:

``` text
ShootEvent
```

The event handler resolves the necessary state:

``` text
Shooter Entity
    ↓
TransformComponent
    ↓
WeaponComponent
```

Suppose:

``` text
Player position:
    (400, 300)

Mouse position:
    (800, 450)
```

The direction is:

``` text
direction =
    normalize(
        target - origin
    )
```

The event processor can then generate:

``` text
SpawnEntityEvent
```

The entity manager processes the spawn request.

The bullet is created.

------------------------------------------------------------------------

# 11. Simulation

Simulation is independent of whether an external command occurred.

Every frame, systems may run:

``` text
AI
Movement
Physics
Animation
Timers
Particles
```

For example:

``` text
AI System:
    Enemy moves toward Player.

Movement System:
    Bullet moves in its velocity direction.
    Enemy moves according to its AI output.

Collision System:
    Detects overlap between entities.
```

The important distinction is:

``` text
Events:
    Discrete changes or requests.

Simulation:
    Continuous evolution of the world.
```

For example:

``` text
ShootEvent
    → creates Bullet

MovementSystem
    → moves Bullet every frame
```

The movement does not need a new event every frame.

------------------------------------------------------------------------

# 12. Collision Processing

Suppose:

``` text
Bullet 100
```

collides with:

``` text
Enemy 20
```

The collision system emits:

```
CollisionDetectedEvent 
```

The collision system should not necessarily know:

``` text
Bullet + Enemy = destroy enemy
```

Instead, the event can be processed using the components on the
entities.

Example:

``` text
Bullet 100:
    DamageComponent

Enemy 20:
    HealthComponent
```

The system can derive:

``` text
DamageEvent
```

and potentially:

``` text
DestroyEntityEvent
```

For example:

``` text
CollisionDetected
        ↓
Component inspection
        ↓
DamageComponent + HealthComponent
        ↓
DamageEvent
```

Then:

``` text
DamageEvent
        ↓
HealthComponent
        ↓
Health decreases
```

If health reaches zero:

``` text
EntityDiedEvent
```

may be emitted.

The bullet can separately produce:

``` text
DestroyEntityEvent
```

------------------------------------------------------------------------

# 13. Event Processing and Deferred Events

Some events may not be processed immediately.

For example:

``` text
CollisionDetected
```

may occur after the primary event phase has already completed.

The event can be queued for the next frame.

A simple model:

``` text
Current Event Queue
    ↓
Events processed this frame

Deferred Event Queue
    ↓
Events processed next frame
```

Example:

``` text
Frame N:

Simulation
    ↓
Collision detected
    ↓
CollisionDetectedEvent
    ↓
Deferred queue
```

Then:

``` text
Frame N + 1:

Event processing
    ↓
CollisionDetectedEvent
    ↓
DamageEvent
    ↓
DestroyEntityEvent
```

This avoids forcing every system to process all events immediately.

------------------------------------------------------------------------

# 14. Full Example Lifecycle

``` text
FRAME N
────────────────────────────────────────────

1. UI INPUT
    LEFT CLICK at (800, 450)

2. COMMAND REQUEST
    CommandRequest:
        POINTER_BUTTON_PRESSED
        Button = LEFT
        Position = (800, 450)

3. QUEUE
    CommandRequestQueue

4. COMMAND PROCESSING
    Game Loop dequeues request.

5. INTERPRETATION
    Interpreter:
        Scan entities.
        Player:
            Transform
            Weapon ← Interactable
            Health
        Weapon matches input.
        Resolve:
            Source = Player
            Position = Mouse position
        Weapon generates:
            ShootEvent

6. EVENT PROCESSING
    ShootEvent:
        Resolve player position.
        Calculate direction.
        Generate:
            SpawnEntityEvent

7. SPAWN
    Bullet is created.
    Bullet:
        Transform
        Movement
        Collision
        Damage

8. SIMULATION
    AI:
        Enemy moves toward player.
    Movement:
        Bullet moves forward.
        Enemy moves.
    Collision:
        Bullet intersects Enemy.

9. DETECTION EVENT
    CollisionDetectedEvent:
        Bullet 100
        Enemy 20

10. EVENT QUEUE
    If current event phase is exhausted:
        Queue for next frame.

11. NEXT FRAME
    Collision event processed.
    Components are inspected:
        Bullet:
            DamageComponent
        Enemy:
            HealthComponent
    Generate:
        DamageEvent
        DestroyEntityEvent

12. STATE UPDATE
    Enemy health decreases.
    Bullet is removed.

13. RENDER SNAPSHOT
    Updated state sent to UI.
```

------------------------------------------------------------------------

# 15. Context Resolution

The Interpreter should resolve context in stages.

## Stage 1: Command Context

Information directly from the input:

``` text
Button:
    LEFT
Screen position:
    (800, 450)
Timestamp:
    T
```

## Stage 2: World Context

Information resolved from the current world:

``` text
World position:
    (1200, 650)
Entity under cursor:
    Enemy 20
Active entity:
    Player 42
```

## Stage 3: Component Context

Information relevant to the interacting component:

``` text
WeaponComponent:
    Owner:
        Player 42
Origin:
    Player.Transform.Position
Target:
    Enemy 20
```

## Stage 4: Event Context

The component produces the smallest useful event:

``` text
ShootEvent {
    ShooterID: 42
    Target: Enemy 20
}
```

The general flow is:

``` text
Command Context
    ↓
World Context
    ↓
Component Context
    ↓
Event Context
```

Context should become progressively more specific.

The entire game state should not be injected into every event.

------------------------------------------------------------------------

# 16. The Builder Game Loop

A possible game loop:

``` text
┌─────────────────────────────┐
│ Input Collection            │
└─────────────┬───────────────┘
              ▼
┌─────────────────────────────┐
│ Command Interpretation      │
└─────────────┬───────────────┘
              ▼
┌─────────────────────────────┐
│ Discrete Event Processing   │
└─────────────┬───────────────┘
              ▼
┌─────────────────────────────┐
│ Continuous Simulation       │
└─────────────┬───────────────┘
              ▼
┌─────────────────────────────┐
│ Detection                   │
│ Collision / Trigger / Query │
└─────────────┬───────────────┘
              ▼
┌─────────────────────────────┐
│ Deferred Event Processing   │
└─────────────┬───────────────┘
              ▼
┌─────────────────────────────┐
│ State Finalization          │
└─────────────┬───────────────┘
              ▼
┌─────────────────────────────┐
│ Render Snapshot             │
└─────────────────────────────┘
```

------------------------------------------------------------------------

# 17. State Mutation Policy

A useful rule is:
> Components and systems should not arbitrarily mutate unrelated entities.

Instead, state changes should be performed through:

-   Component-owned methods.
-   Engine operations.
-   State-changing events.
-   Dedicated systems.

For example:

``` text
ShootEvent
    ↓
Weapon Handler
    ↓
SpawnEntityEvent
    ↓
Entity Manager
    ↓
Bullet created
```

Rather than:

``` text
WeaponComponent
    directly modifies global entity map
```

This keeps ownership clear.

------------------------------------------------------------------------

# 18. Example: Same Builder, Different Games

## Top-Down Shooter

``` text
Player
├── Transform
├── Weapon
└── Health

Enemy
├── Transform
├── AI
└── Health
```

Input:

``` text
LEFT_CLICK
```

Result:

``` text
WeaponComponent
    ↓
ShootEvent
```

## Door Puzzle Game

``` text
Player
├── Transform
└── InteractionComponent

Door
├── Transform
└── DoorComponent
```

Input:

``` text
E
```

Result:

``` text
DoorComponent
    ↓
OpenDoorEvent
```

------------------------------------------------------------------------

## Pickup Game

``` text
Player
├── Transform
└── InventoryComponent

Coin
├── Transform
└── PickupComponent
```

Input:

``` text
E
```

Result:

``` text
PickupComponent
    ↓
PickupEvent
```

The Builder remains unchanged.

Only the entities, components, and metadata differ.

------------------------------------------------------------------------

# 19. Key Design Principles

## Principle 1: The UI Sends Facts, Not Game Semantics

Good:

``` text
LEFT_CLICK at (800, 450)
```

Avoid:

``` text
Shoot at Enemy 20
```

## Principle 2: Input Mapping Must Be Replaceable

Changing:

``` text
LEFT_CLICK
```

to:

``` text
RIGHT_CLICK
```

must not require changing the Weapon's task implementation.

## Principle 3: The Interpreter Should Know Interfaces, Not Game Types

The Interpreter knows:

``` text
Interactable
Event
Entity
Component
Context
```

It should not know:

``` text
Weapon
Enemy
Door
Bullet
```

## Principle 4: Components Own Task Semantics

If a component is responsible for a task, it should produce the event
representing that task.

``` text
WeaponComponent
    → ShootEvent
```

``` text
DoorComponent
    → OpenDoorEvent
```

The generic Interpreter should not construct these event types directly.

## Principle 5: Context Is Resolved Progressively

``` text
Command
    ↓
World Context
    ↓
Component Context
    ↓
Event Context
```

Do not inject the entire game state into every event.

## Principle 6: Simulation Is Separate from Discrete Events

``` text
Events:
    Shoot
    Spawn
    Damage
    Destroy

Simulation:
    AI
    Movement
    Physics
    Timers
```

Simulation runs even when no external command occurs.

------------------------------------------------------------------------

# 20. Final Architecture

The complete conceptual architecture is:

``` text
                         ┌────────────────────┐
                         │        UI          │
                         └─────────┬──────────┘
                                   │
                                   ▼
                         ┌────────────────────┐
                         │  CommandRequest    │
                         │                    │
                         │  Raw external      │
                         │  input facts       │
                         └─────────┬──────────┘
                                   │
                                   ▼
                         ┌────────────────────┐
                         │ Command Queue      │
                         └─────────┬──────────┘
                                   │
                                   ▼
                         ┌────────────────────┐
                         │   Interpreter      │
                         │                    │
                         │ Scan candidates    │
                         │ Resolve context    │
                         │ Match interactors  │
                         └─────────┬──────────┘
                                   │
                                   ▼
                         ┌────────────────────┐
                         │ Interactable       │
                         │ Components         │
                         │                    │
                         │ Component-specific │
                         │ task semantics     │
                         └─────────┬──────────┘
                                   │
                                   ▼
                         ┌────────────────────┐
                         │       Events       │
                         └─────────┬──────────┘
                                   │
                                   ▼
                         ┌────────────────────┐
                         │ Event Processor    │
                         └─────────┬──────────┘
                                   │
                                   ▼
                         ┌────────────────────┐
                         │ State Operations   │
                         │                    │
                         │ Spawn              │
                         │ Destroy            │
                         │ Damage             │
                         │ Modify             │
                         └─────────┬──────────┘
                                   │
                                   ▼
                         ┌────────────────────┐
                         │    Simulation      │
                         │                    │
                         │ AI                 │
                         │ Movement           │
                         │ Physics            │
                         │ Timers             │
                         └─────────┬──────────┘
                                   │
                                   ▼
                         ┌────────────────────┐
                         │ Detection Systems  │
                         │                    │
                         │ Collision          │
                         │ Triggers           │
                         └─────────┬──────────┘
                                   │
                                   ▼
                         ┌────────────────────┐
                         │ Deferred Events    │
                         └─────────┬──────────┘
                                   │
                                   ▼
                         ┌────────────────────┐
                         │ Updated Game State │
                         └─────────┬──────────┘
                                   │
                                   ▼
                         ┌────────────────────┐
                         │ Render Snapshot    │
                         └────────────────────┘
```

The central abstraction is:

``` text
External Command
    ↓
Generic Interpreter
    ↓
Interactable Component
    ↓
Task Event
    ↓
Event Processing
    ↓
Simulation
    ↓
State
```

The Builder is therefore not a game engine that understands games.

It is a generic runtime that understands:

``` text
Input
    → Component capability
    → Event
    → State transition
    → Continuous simulation
```

The actual game emerges from the component composition and metadata
created by the user.
