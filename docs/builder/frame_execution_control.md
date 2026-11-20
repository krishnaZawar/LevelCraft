# Runtime Execution Control — Design Document

## 1. Context and Motivation

The game loop currently starts as part of backend initialization, independently of whether the frontend has established its WebSocket connection or finished initializing the runtime view.

This creates a timing problem: the backend can begin processing game frames before the user is ready to interact with the game. Simulations, movement, and other continuous game behaviours may advance while the frontend is still connecting or initializing.

Runtime execution therefore needs an explicit readiness gate.

The objective is to ensure that:

- The backend can initialize the game and start its game loop without immediately advancing gameplay.
- Gameplay processing begins only after the frontend signals that it is ready.
- The backend can continue receiving and processing runtime-control requests while gameplay is inactive.
- Integration tests can execute the game loop without requiring a frontend or WebSocket connection.

The solution introduces a runtime execution-control mechanism that determines whether the game loop is permitted to process gameplay for a frame.

## 2. Design Overview

The design separates runtime execution control from gameplay execution.

There are two primary responsibilities:

**Control Manager**
- Owns the current runtime execution state.
- Determines whether gameplay processing is permitted.
- Provides a central interface for changing and querying execution state.

**Game Loop**
- Owns the actual execution of game frames.
- Checks the execution state before processing gameplay.
- Continues handling the operations necessary to manage runtime execution, even when gameplay is inactive.

The control manager does not execute, schedule, or manage individual gameplay systems. It only provides the decision that determines whether gameplay processing is allowed.

The game loop remains the authority responsible for executing gameplay.

### 2.1 Execution State

The initial implementation uses two states:

| State | Meaning |
|---|---|
| Inactive | Gameplay processing is suspended. |
| Active | Gameplay processing is permitted. |

The default state is **Inactive**.

This ensures that starting the backend does not implicitly start gameplay. Execution begins only after an explicit activation request.

## 3. Low-Level Design

### 3.1 Control Manager

The Control Manager is a runtime-level subsystem responsible for maintaining execution state.

Its responsibilities are deliberately limited:

- Maintain the current execution state.
- Validate state transitions or state updates.
- Expose the current state.
- Provide a simple check to determine whether a gameplay frame can execute.

The manager should own its state privately and expose controlled methods for reading and updating it.

Its state is accessed and modified by the game-loop goroutine. Runtime requests received through the WebSocket are placed into the command queue and processed by that same goroutine. Consequently, reads and writes to the manager are serialized.

This single-goroutine ownership is an architectural invariant. If other goroutines begin accessing or modifying the manager directly, synchronization or a different ownership model will be necessary.

### 3.2 Execution-State Command

Runtime execution state is changed through the existing command-processing mechanism.

A dedicated execution-state command is responsible for:

1. Receiving the requested execution state.
2. Validating the request.
3. Updating the Control Manager.
4. Returning an error if the requested state is invalid.

The WebSocket handler does not modify the Control Manager directly. It only receives the request and places it into the command queue.

This preserves the existing separation between transport handling and runtime state management.

It also ensures that execution-state changes are applied in the same execution context that evaluates the gate.

### 3.3 Integration with the Game Loop

The execution gate belongs inside the game loop, after runtime-control commands have been processed and before gameplay processing begins.

The conceptual frame flow becomes:

1. Process incoming commands.
2. Check the current execution state.
3. If inactive:
   - Reset input states.
   - Skip gameplay processing.
4. If active:
   - Generate input-based events.
   - Process generated events.
   - Execute continuous simulations.
   - Process simulation-generated events.
   - Detect and resolve collisions.
5. Publish the current game-state snapshot.
6. Maintain frame pacing and proceed to the next iteration.

The gate controls gameplay processing rather than the entire game loop.

This distinction is important because the game loop must continue running while gameplay is inactive. Otherwise, it would be unable to process the command that activates execution again.

### 3.4 Command Processing While Inactive

Control commands must remain processable while the game is inactive. This is necessary for activation and other future runtime-control operations.

However, allowing every command to execute before the gate introduces an important limitation: gameplay-related commands can also be consumed while execution is inactive.

The current design should therefore distinguish between two conceptual categories:

**Runtime-control commands**
- Change execution state.
- Manage runtime-level behaviour.
- Must remain processable while gameplay is inactive.

**Gameplay commands**
- Affect gameplay state or trigger gameplay behaviour.
- Should respect the execution gate.

As the command system expands, this distinction should be reflected in command-processing policy. A separate control-command processing path or an explicit command classification can achieve this.

The gate should not simply be moved before all command processing, because doing so would prevent the inactive game loop from receiving its own activation request.

### 3.5 Game-State Snapshot Publication

Snapshot publication remains outside the gameplay execution gate.

Even while gameplay is inactive, the backend must be able to publish the initial game state to the frontend.

This allows the frontend to initialize and render the scene before gameplay begins.

Keeping snapshot publication outside the gate also means the frontend can continue receiving state updates independently of whether gameplay processing is active.

### 3.6 Delta-Time Handling

Frame timing must continue to be updated while execution is inactive.

The time between frames should not accumulate into the first active frame. Otherwise, resuming gameplay could produce an unusually large delta time, causing simulations and movement to advance unexpectedly.

The game loop should therefore continue maintaining its frame timing while inactive, ensuring that the first active frame uses a normal frame delta.

### 3.7 Input-State Handling

Input states should be reset while gameplay is inactive.

Without this reset, a key pressed before deactivation could remain registered as pressed when execution resumes.

Resetting input states prevents stale input from affecting the first active frames.

## 4. End-to-End Execution Flow

The intended lifecycle is as follows:

### Backend initialization

1. Load and validate the game project.
2. Initialize the runtime systems and Control Manager.
3. Initialize the execution state as Inactive.
4. Start the game loop.
5. Start the backend server.

The game loop can begin running immediately, but gameplay remains suspended.

### Frontend initialization

1. The frontend establishes a WebSocket connection.
2. The backend accepts the connection.
3. The frontend receives the initial game-state snapshot.
4. The frontend initializes the runtime view and completes the required rendering setup.
5. Once the frontend is ready, it sends an execution-state activation request.

The activation request is processed through the normal command pipeline.

An important distinction is that a successful WebSocket connection alone does not establish frontend readiness. Transport connectivity and runtime initialization are separate milestones. Activation should represent completion of the latter.

### Gameplay execution

1. The game loop processes the activation command.
2. The Control Manager transitions to Active.
3. The gate permits gameplay processing.
4. Subsequent frames execute the normal gameplay pipeline.

The backend remains responsible for executing frames; the frontend only communicates when execution is permitted to begin.

## 5. Integration Testing

The execution gate should not introduce a dependency on the frontend or WebSocket connection.

Integration tests should be able to:

1. Initialize the game state and runtime systems.
2. Set the Control Manager to Active.
3. Run the game loop or the relevant frame-processing path.
4. Verify gameplay behaviour and resulting state changes.

Tests that need to verify inactive behaviour can leave the manager in its default Inactive state and assert that gameplay processing does not occur.

This keeps the runtime execution mechanism independent of the UI and allows backend integration tests to remain self-contained.

## 6. Future Extensibility

The initial design intentionally supports only Active and Inactive states. However, centralizing execution control creates a natural extension point for future runtime-management capabilities.

Potential extensions include:

### Pause and Resume

Pause and resume can use the existing execution-state mechanism without requiring individual gameplay systems to implement their own pause logic.

### Debug Execution

A future execution mode could distinguish normal execution from debug execution.

For example, debug execution might allow additional inspection or controlled frame progression. Execution state and execution mode should remain separate concepts: one determines whether gameplay can proceed, while the other determines how it proceeds.

### Single-Frame Execution

A future frame-stepping capability could allow the user to advance gameplay one frame at a time.

This would require a consumable frame-step request or permit rather than simply introducing another boolean state. A boolean could collapse multiple step requests into one.

### Runtime Inspection

The control mechanism could later coordinate with runtime inspection capabilities, allowing the frontend to inspect game state while gameplay is paused.

### Connection Lifecycle

The current design activates execution based on an explicit frontend readiness request. A future lifecycle extension could also handle disconnection and reconnection.

If the intended behaviour is to pause gameplay when the frontend disconnects, that policy should be handled explicitly rather than assuming that connection loss automatically changes execution state.

On reconnection, the frontend should receive a fresh snapshot and complete initialization before requesting activation again.

These capabilities can build on the existing control boundary without moving gameplay responsibilities into the Control Manager.

## 7. Design Boundaries

The design follows a few important ownership principles:

- **Control Manager:** Owns runtime execution state and answers whether gameplay is permitted.
- **Game Loop:** Owns frame execution and applies the execution gate.
- **Command Processing:** Provides the mechanism through which runtime-control requests reach the manager.
- **WebSocket Handler:** Handles transport and request ingestion; it does not directly control gameplay execution.
- **Frontend:** Determines when its runtime is ready and communicates that readiness to the backend.

The Control Manager should not become a general-purpose runtime coordinator. Its responsibility is execution control, not managing simulations, input, collisions, rendering, or individual components.

Similarly, gameplay systems should not independently maintain their own copies of the global execution state. The game loop provides the central execution boundary.

## 8. Summary

The runtime execution gate prevents gameplay from advancing before the frontend is ready.

The design achieves this by introducing a dedicated Control Manager with an initially inactive execution state. The game loop continues running, processes runtime-control requests, and checks the manager before entering the gameplay pipeline.

Keeping command processing, gameplay execution, and snapshot publication conceptually separate allows the backend to initialize and communicate with the frontend without prematurely advancing the game.

Most importantly, the mechanism remains independent of the UI. This supports backend-only integration testing today and provides a clear foundation for future runtime capabilities such as pause/resume, debug execution, frame stepping, and runtime inspection.