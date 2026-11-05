# LevelCraft — Collision System Design

## 1. Overview

The Collision System is responsible for detecting and resolving collisions between game objects.

It uses **Axis-Aligned Bounding Box (AABB)** geometry. Collider participation and movement eligibility are determined by the Collider component, while **collision boundaries are derived from Transform**.

The implementation focuses on simple positional correction rather than sophisticated physics simulation.

## 2. Collider Design

- Collider determines whether an object participates in collision detection.
- Transform provides the object's position and dimensions.
- Collider boundaries are derived entirely from Transform; users do not configure them independently.
- An object without either Collider or Transform is excluded from collision detection.
- Objects currently use square, axis-aligned boundaries.

## 3. Game Loop Integration

Collision detection and resolution happen after:

1. CommandRequest execution.
2. Event execution.
3. Simulation execution.

The collision phase then performs detection followed by resolution.

## 4. Collision Detection

**The system examines every unique pair of eligible game objects using an all-pairs approach.**

AABB geometry is used to calculate four directional penetration values: left, right, top and bottom.

The horizontal and vertical intersection depths are derived from these values. A collision is registered only when both are strictly positive. **Edge-touching is not considered a collision**.

Each detected pair generates a `CollisionOccurred` event.

Detection operates on a snapshot of the game state at the beginning of the collision phase.

## 5. Collision Event Processing

Detected collision events are placed into two queues:

- **Dedicated collision queue:** Prioritizes collision resolution and is processed completely during the collision phase.
- **Standard event queue:** Allows other components subscribed to collision events to react, subject to the normal event-processing limit.

> **Note**: The same collision must not be resolved twice through these two processing paths.

## 6. Collision Resolution

For every collision event, the handler:

1. Retrieves the participating objects and their components.
2. Recalculates overlaps using their current Transform values.
3. Skips resolution if the objects no longer overlap.
4. Determines the minimum positional correction using AABB geometry.
5. Distributes and applies the correction according to movement eligibility.

> **Note**: This distribution and transform update happen in the same step. The handler does not emit another event for the transform to update because that could lead to stale calculations and wrong updates

### Movement Rules

| Object A | Object B | Resolution |
|---|---|---|
| Movable | Movable | Equal correction |
| Movable | Immovable | A receives full correction |
| Immovable | Movable | B receives full correction |
| Immovable | Immovable | No correction |

Transform updates are applied synchronously by the resolution handler. **Positional corrections are geometric displacements and are not scaled by delta time because they are precisely for that frame only and not to be smoothened out.**

The first object is used as the reference when calculating correction directions. This does not imply priority or unequal treatment.

## 7. Multiple Collisions

**Collision detection is snapshot-based, but resolution is sequential.**

Each resolution uses the latest Transform values. This allows earlier corrections to eliminate subsequent collisions without applying unnecessary updates.

A correction can introduce a new collision with a pair that was not detected initially. Such collisions are discovered during a subsequent detection pass.

> **Note**: The initial implementation does not guarantee that every collision in the scene will be resolved in a single pass.

## 8. Scope and Limitations

The initial Collision System supports:

- Square AABB collision detection.
- Positional collision resolution.
- Movable and immovable objects.
- Multiple collision pairs within a frame.
- Event-driven collision notifications.

The following are outside the current scope:

- Velocity and momentum-based responses.
- Friction and bouncing.
- Continuous collision detection.
- Spatial partitioning.
- Iterative collision solving.
- Collision lifecycle tracking (started, continued, ended).

These can be considered as the engine's requirements evolve.

## 9. Design Principles

The implementation should preserve the following principles:

- Transform remains the single source of truth for collider boundaries.
- Detection remains independent of the mechanism that moved an object.
- Every unique pair is detected once per pass.
- Resolution uses current state and applies corrections synchronously.
- Collision resolution remains separate from other components' reactions to collision events.
