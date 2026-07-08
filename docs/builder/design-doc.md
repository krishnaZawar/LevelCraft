# Game Builder

This document describes the architecture of a command-driven, event-based game engine designed for deterministic state updates, modular components, and clear separation between input, game logic, and simulation.

The core philosophy is:

> The UI expresses intent. The backend determines meaning. Components own state. Events connect components.

The architecture separates:

* **Commands** — requests representing intent.
* **Events** — facts representing state changes or cross-component communication.
* **Components** — independent systems responsible for owning and updating their own state.
* **Game Loop** — the authoritative executor of all game state changes.

## Design Goals
### 1. Deterministic Execution

The engine should process actions in a predictable order.

A command must be fully resolved before another command begins:
```
Command A
    |
    v
Events
    |
    v
State Updates
    |
    v
Event Chain Complete

Command B
```
This avoids ordering issues where multiple actions partially affect the game state.

### 2. Component Independence

Components should not directly modify each other's state.

Example:

- Incorrect:
    ```
    AttackComponent
            |
            v
    HealthComponent.health -= damage
    ```
- Correct:
    ```
    AttackComponent
            |
            v
    DamageAppliedEvent
            |
            v
    HealthComponent
            |
            v
    Update Health
    ```
Each component owns its own data and exposes behavior through events.

### 3. Generic UI Layer

The UI should not contain game-specific knowledge.

The UI only generates generic input commands.

Example:
```
Player presses button
        |
        v
Generic Input Command
        |
        v
Backend Interpretation
        |
        v
Game-Specific Events
```
The same UI architecture can support different games because the backend determines what an input means.

## Threading Model

The engine uses two conceptual areas:

### 1. Main Thread

Responsibilities:

- Receive UI commands.
- Insert commands into the command queue.
- Handle UI state requests.
- Return state snapshots.

The main thread does not:

- Run gameplay logic.
- Modify game state.
- Execute events.


### 2. Game Loop Thread

The game loop owns the authoritative game state.

Responsibilities:

* Consume commands.
* Validate commands.
* Generate events.
* Process event chains.
* Run simulations.
* Update components.

Only the game loop may mutate game state.


## Frame Execution Lifecycle

Each frame follows this order:

```
Frame Start
    |
    v
Command Processing Phase
    |
    +--> Consume commands
    |
    +--> Validate command
    |
    +--> Generate events
    |
    +--> Drain event queue completely
    |
    v
Simulation Phase
    |
    +--> Physics
    |
    +--> AI
    |
    +--> Movement
    |
    +--> Timers
    |
    v
Simulation Event Resolution
    |
    +--> Process generated events
    |
    v
Publish State Snapshot
    |
    v
Frame End
```

> Note: The simulation event resolution can be done at the same time or once all the events are published. The decision is yet to made

## Command System

### Command Flow

```
UI Input
    |
    v
Generic Command
    |
    v
Command Queue
    |
    v
Game Loop
    |
    v
Validation
    |
    v
Events
```

### Command Processing Rules

#### **Sequential Processing**

Commands are processed sequentially.

A command is not considered complete until:

- Validation is complete.
- All generated events are processed.
- All resulting event chains are resolved.

#### **Command Limit**

The initial implementation uses a configurable command limit per frame.

Example:

```
Maximum Commands Per Frame = N
```

Future improvement:

```
Process commands until time budget is exhausted
```


### Invalid Commands

Invalid commands are discarded.

Example:

```
AttackCommand

Target does not exist
        |
        v
    Discard
```

Invalid commands:

- Do not emit events.
- Do not modify state.
- Do not affect gameplay flow.


## Event System

### Event Queue Ownership

The game loop owns a single global event queue.

Components do not have independent queues.

```
Component
    |
    v
Event Queue
    |
    v
Event Handler
    |
    v
Component Update
```

## Event Processing Rules

When a command generates events:

1. Add events to the event queue.
2. Process events sequentially.
3. Allow event handlers to generate additional events.
4. Continue until the queue is empty.

Example:

```
AttackCommand
      |
      v
EnemyAttackEvent
      |
      v
AttackComponent
      |
      v
DamageAppliedEvent
      |
      v
HealthComponent
      |
      v
DestroyGameObjectEvent
      |
      v
GameObjectManager
```

Only after the queue is empty can the next command execute.


## Component Model

Each component owns its own state.

Example:

```
HealthComponent

Owns:
- Current HP
- Max HP
- Damage rules
```

Other components cannot directly modify health.

They communicate through events.


## Event Emission Rules

Events exist only for cross-component communication.

A component should emit an event when another component needs to react.

Example:

Movement:

```
Position += velocity * deltaTime
```

No event required.


Collision:

```
Collision detected
        |
        v
CollisionEvent
```

Event required because another component may need to react.


> Note: Internal updates stay internal. Cross-component updates become events.


## Simulation Model

Simulation systems update internally.

Examples:

- Physics
- AI
- Movement
- Timers
- Animation

These systems do not generate events for every update.

Bad:

```
PositionChangedEvent
PositionChangedEvent
PositionChangedEvent
```

for every frame.

Good:

```
CollisionDetectedEvent
EntityDestroyedEvent
QuestCompletedEvent
```

Only meaningful cross-component changes enter the event system.


## UI Communication

### Initial Approach: Pull-Based Snapshots

The UI requests state snapshots.

```
 UI
 |
 | State Request
 v
Main Thread
 |
 v
State Snapshot
```

The UI does not directly access game state.


### Future Option: Separate Communication Channels

Possible future architecture:

Command channel:

```
 UI
 |
 | Commands
 v
WebSocket
 |
 v
Backend
```

State channel:

```
Backend
 |
 | State Updates
 v
SSE / Dedicated Stream
 |
 v
UI
```

This allows independent evolution of command and state synchronization.

## Example Flow

### Enemy Attack Example

#### **Step 1**

UI:

```
Player clicks attack button
```

Creates:

```
Generic Input Command
```

#### **Step 2**

Backend interprets command:

```
Input Command
        |
        v
EnemyAttackEvent
```


#### **Step 3**

Attack component handles event:

```
EnemyAttackEvent
        |
        v
AttackComponent
```

Creates:

```
DamageAppliedEvent
```


#### **Step 4**

Health component handles damage:

```
DamageAppliedEvent
        |
        v
HealthComponent
```

Updates:

```
Health = 0
```

Creates:

```
DestroyGameObjectEvent
```


#### **Step 5**

Game object manager handles destruction:

```
DestroyGameObjectEvent
        |
        v
GameObjectManager
        |
        v
Remove Entity
```

The command is now fully resolved.


## Benefits

- Clear Ownership
    - Each component has a defined responsibility.
- Reduced Coupling
    - Components communicate without knowing implementation details of other components.
- Deterministic Ordering
    - Commands and events execute in a controlled order.
- Debugging
    - The engine can inspect:

        ```
        Command
            |
            v
        Events
            |
            v
        State Changes
        ```
        making issues easier to trace.
- Extensibility
    - Future additions can be introduced without changing existing systems:
        - AI
        - Networking
        - Replay
        - Multiplayer
        - Modding
        - New game rules

## Future Extensions

The following are intentionally left open:

- Time-Based Command Scheduling
    - Replace command count limits with CPU time budgets.
- Parallel Event Handling
    - Multiple subscribers for the same event may be introduced later.
- Replay System
    - The deterministic command/event flow can support recording and replay.
- Network Replication
    - The separation between commands and state snapshots allows future multiplayer synchronization.
