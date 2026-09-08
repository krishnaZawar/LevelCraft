# Movement System Design

## 1. Overview

Movement in LevelCraft should be modelled as a **composable capability** rather than as predefined movement types such as `TopDownMovement`, `PlatformerMovement`, `FollowMovement`, etc.

An entity's movement behaviour should emerge from the combination of components attached to it.

The model needs to answer three fundamental questions:

1. **Can the entity move?**
2. **What determines its movement?**
3. **What can influence or constrain that movement?**

The engine should ultimately resolve these inputs into the entity's actual movement.

## 2. Movement Model

A `Movement` component represents an entity's ability to move and contains the state required to perform movement.

Conceptually:

```text
Movement
├── velocity
├── speed
├── acceleration
└── movement configuration
```

Other components interact with this state rather than implementing their own movement mechanics.

There are two broad categories of components involved:

### Movement Producers

Determine where the entity wants to move.

Examples:

```text
InputMovement
Follow
Patrol
MoveTo
```

### Movement Modifiers

Modify movement independently of where the entity wants to go.

Examples:

```text
Gravity
Knockback
Friction
```

This distinction allows different behaviours to be composed without creating movement-specific components.

## 3. Input-Driven Movement

Input-dependent movement is handled through the existing command/event mechanism.

An input command is matched against the input configuration of relevant components. If a component responds to that input, it generates an appropriate event.

For example:

```text
Input Command
      ↓
InputMovement component matches input
      ↓
Move event
      ↓
Movement state updated
```

The event should contain the movement information required by the component, such as direction or an action such as jump.

The `InputMovement` component therefore owns the relationship between input and movement behaviour.

It should not need to know about the specific game type.

## 4. Continuous Movement

Movement that does not depend on an explicit input command is handled during continuous simulation.

Examples include:

```text
Follow → determine direction toward target
Patrol → determine direction toward next waypoint
Gravity → continuously modify velocity
```

These components operate on the entity's existing movement state.

For example:

```text
Follow
  ↓
calculate desired direction
  ↓
update Movement
```

and:

```text
Gravity
  ↓
apply acceleration
  ↓
update Movement.velocity
```

They do not need to generate input events because their behaviour exists independently of user input.

## 5. Movement Resolution

The `Movement` component provides the common state through which different movement behaviours interact.

A simplified representation is:

```text
Movement
    ↑
    │
 ┌──┴───────────────┐
 │                  │
InputMovement    Follow/Patrol
 │                  │
 └───────┬──────────┘
         │
         ▼
      velocity
         ↑
         │
   Gravity/Forces
```

The resulting velocity/displacement represents what the entity is attempting to do.

Collision handling subsequently determines how much of that movement can actually be applied.

```text
Movement
    ↓
Proposed displacement
    ↓
Collision resolution
    ↓
Actual displacement
```

Movement components therefore do not need to perform collision detection themselves.

## 6. Composition Examples

### Top-down player

```text
Entity
├── Transform
├── Movement
└── InputMovement
```

Input determines the movement direction.

### Platformer player

```text
Entity
├── Transform
├── Movement
├── InputMovement
├── Gravity
└── Jump
```

Horizontal movement comes from input, gravity continuously affects velocity, and jump modifies vertical movement.

There is no `PlatformerMovement` component.

### Following enemy

```text
Entity
├── Transform
├── Movement
└── Follow
```

`Follow` continuously determines the desired direction toward its target.

Additional behaviour can be composed naturally:

```text
Entity
├── Transform
├── Movement
├── Follow
├── Gravity
└── Collider
```

## 7. Position in the Existing Game Loop

Movement participates in the existing loop in two places:

```text
Command
  ↓
Generate Events
  ↓
InputMovement events
  ↓
Execute Events
  ↓
Continuous Simulations
  ├── Follow / Patrol
  ├── Gravity
  └── other movement influences
  ↓
Collision Resolution
  ↓
Final movement / transform update
```

The game loop itself remains unchanged. Movement simply provides components that participate in the appropriate existing stages.

## 8. Design Principles

### No predefined movement types

Avoid:

```text
MovementType = TOP_DOWN
MovementType = PLATFORMER
MovementType = FOLLOW
```

These behaviours should emerge through component composition.

### Components should own their behaviour

`InputMovement` owns input-to-movement behaviour.

`Follow` owns target-following behaviour.

`Gravity` owns gravitational influence.

They should not become a collection of special cases inside a central movement component.

### Movement remains the shared state

Different movement-related components should converge on the entity's `Movement` state rather than maintaining independent movement implementations.

### Collision remains independent

Movement determines the proposed motion.

Collision resolution determines whether that motion can occur.

### Avoid premature abstractions

Concepts such as *intent*, *forces*, *constraints*, and *movement solvers* are useful for reasoning about the system, but should not automatically become separate engine abstractions.

Introduce additional layers only when the implementation demonstrates a concrete need for them.

## 9. Result

The resulting architecture allows LevelCraft to express movement through composition:

```text
Movement + InputMovement
        → Top-down movement

Movement + InputMovement + Gravity + Jump
        → Platformer movement

Movement + Follow
        → Following movement

Movement + Patrol
        → Patrol movement
```

This keeps movement extensible while fitting naturally into LevelCraft's existing entity/component and event/simulation architecture.
