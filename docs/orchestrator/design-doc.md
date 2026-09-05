# LevelCraft Orchestrator

The orchestrator is a standalone service that is used to coordinate and manage the process lifecycles of independent services.

## Why is the orchestrator needed?

Our application being a multi-process application would not be as simple to start and manage as a single process application. Things that add in to the management of multi process applications are:
- We cannot run all the individual processes manually and coordinate everytime we want to run an application. This would be cumbersome and user unfriendly
- Managing the order of startup of processes
    - During the startup of an application, the processes should be created in order for proper linkage and startup
    - Addition of new components would be easier to manage via the orchestrator
- Cleanup would be easier to manage
    - Whenever we want to close the application, or parts of the application crash, We can ensure the other running processes are exited properly and the application does not leave out memory or process leaks and we achieve a graceful shutdown
- Allows flexible handling of specific components
    - Say a process running is not very important, so even if that process corrupts or crashes, it should not crash the application. In such cases we can explicitly handle that component in that way
    - Same way, say we have a process that if crashed, the application should shut down, we can handle it that way respectively
    - We can have robust process handling, restarts, etc. based on the need
- Any process creation inside the application should go through the orchestrator, in order to maintain the record of all the spawned processes
- We can have healthchecks in place to track process liveness and handle things accordingly.

## The 3 flows that define the Orchestrator

The orchestrator can be broadly divided into three flows:

1. **Startup Flow** – Starts and initializes all required processes in the correct order.
2. **Runtime Flow** – Monitors and manages processes while the application is running.
3. **Failure Handling / Exit Flow** – Handles process failures and coordinates application shutdown.

```text
                    ┌──────────────────────┐
                    │      Application     │
                    └──────────┬───────────┘
                               │
                               ▼
                    ┌──────────────────────┐
                    │     Orchestrator     │
                    └──────────┬───────────┘
                               │
             ┌─────────────────┼─────────────────┐
             │                 │                 │
             ▼                 ▼                 ▼
        Startup Flow      Runtime Flow     Failure / Exit
             │                 │                 │
             ▼                 ▼                 ▼
       Start processes     Monitor          Restart / Stop
       in order            processes        processes
             │                 │                 │
             └─────────────────┼─────────────────┘
                               ▼
                    ┌──────────────────────┐
                    │  Managed Processes   │
                    └──────────────────────┘
```

### 1. Startup Flow

The startup flow starts the required processes in the correct order and waits for each process to be ready before starting dependent processes.

Look at the example below:

```text
Application Start
       │
       ▼
  Orchestrator
       │
       ▼
Build Dependency Graph
       │
       ▼
Determine Startup Order
       │
       ▼
Start Process
       │
       ▼
Wait for Process Ready
       │
       ▼
Start Dependent Process
       │
       ▼
   All Processes
       │
       ▼
Application Ready
```

### 2. Runtime Flow

Once the application is running, the orchestrator monitors all processes it started.

```text
              ┌──────────────┐
              │ Orchestrator │
              └──────┬───────┘
                     │
                  Monitor
                     │
        ┌────────────┼────────────┐
        ▼            ▼            ▼
       UI         Backend       Worker
        │            │            │
        └────────────┼────────────┘
                     │
                     ▼
              Process Healthy?
                │          │
               Yes         No
                │          │
                ▼          ▼
             Continue   Failure Flow
```

The orchestrator monitors:
- Whether processes are alive
- Whether processes are healthy
- Process exits
- Process state

As long as the processes are healthy, the application continues running.

### 3. Failure Handling / Exit Flow

The orchestrator handles both unexpected process failures and intentional application shutdown.

#### failure flow
```text
Process Failure
      │
      ▼
Orchestrator
      │
      ▼
Check Failure Policy
      │
   ┌──┴───────┐
   ▼          ▼
Restart    Shutdown
   │          │
   ▼          ▼
Process     Stop All
Restarts    Processes
   │
   ▼
Process Ready
   │
   ▼
Continue
```

Depending on the process, a failure can result in:
- Restarting the process
- Ignoring the failure
- Marking the application as degraded
- Shutting down the application

#### exit flow

```text
Application Exit
       │
       ▼
 Orchestrator
       │
       ▼
Stop Processes
       │
       ▼
Wait for Exit
       │
       ▼
Force Stop if Needed
       │
       ▼
All Processes Stopped
       │
       ▼
Application Exit
```

### Overall Flow

```text
             ┌──────────────┐
             │    START     │
             └──────┬───────┘
                    ▼
             ┌──────────────┐
             │   STARTUP    │
             └──────┬───────┘
                    ▼
             ┌──────────────┐
             │   RUNTIME    │◄──────────┐
             └──────┬───────┘           │
                    │                   │
              Failure / Exit            │
                    │                   │
                    ▼                   │
             ┌──────────────┐           │
             │   HANDLING   │           │
             └──────┬───────┘           │
                    │                   │
             ┌──────┴──────┐            │
             ▼             ▼            │
          Restart       Shutdown        │
             │             │            │
             ▼             ▼            │
          Runtime       Cleanup         │
             │                          │
             └──────────────────────────┘

```