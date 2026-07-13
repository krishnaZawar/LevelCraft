# CommandRequests for Editor

The editor sends out commands that modify the state of the game. These commands are then validated and the game state is updated accordingly.

These commands should have all the necessary information needed for the update and should be sent out as a request to the backend, where it is broken down into relevant events for update.

The commandRequest should be sent out as a JSON. the structure, for example, would look something like this:
```json
{
    "requestType": "UpdateComponent",
    "requestDetails": {
        "objectID": "obj123",
        "componentName": "shoot",
        "data": {
            // updated component details as key value pairs
        }
    }
}
```

# CommandRequests for Builder

The UI sends out input commandRequests to the builder to validate the request and send out relevant events.

Here the commandRequests vary from how they are sent in the editor. They vary in the context that they carry.
- In the editor the commandRequest will be more specific to the action they are performing in the UI as these actions are deterministic in the editor.
-  For the builder only the generic input commandRequests are sent because they do not know which for which game the commandRequest is sent out and the meaning from the request is extracted during runtime based on the games configs.

The structure here as well would be a JSON. It would like something like this:
- for keyboard events
    ```json
    {
        "requestType": "KeyPressed",
        "requestDetails": {
            "key": "W"
        }
    }
    ```
- for mouse press events
    ```json
    {
        "requestType": "MouseButtonPressed",
        "requestDetails": {
            "button": "left",
            "pos": {
                "x": 120,
                "y": 250
            }
        }
    }
    ```

> Note: The CommandRequests for editor are specific to the actions because the actions are deterministic and defined by us. But for the builder, the CommandRequests are generic because the user can create any time of game and the builder should be able to support all of them. Hence for builder, the meaning to the CommandRequest is added during the runtime.