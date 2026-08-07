# Game Editor Design

Before getting into the editors detailed design, let us first understand what all components exist and what are the tasks of each of these components

On a very high level, the game editor has 2 main components:
- Frontend / UI
    - displays the game scene and all the editor components visually to the user
    - allows selection of existing gameObjects in the editor
    - implements drag and drop functionality for the user
    - propagate state update requests to the backend for validation and changes
- Backend
    - accepts and processes command requests from the frontend
    - validates/updates any changes done to the game scene

## Frontend

The panel updates and all the UI are handled in the frontend.
These updates include:
- Selection/Deselection of gameObjects in one panel and updating the others respectively.
- drag and drop of objects in the workspace panel.
- Command Request dispatch according to the event occured
    - These requests wont be sent for every small operation, it will be sent only on notable events that change the game scene state

> Note: For the drag and drop, once the operation is complete there should be a commandRequest propagated to update the new state of the object.

## Backend

The backend should not worry about the styling and placement of the panels nor about the selection/deselection or drag and drop events performed on the frontend. The backend should mainly listen for the following things:
- CommandRequests from the frontend to process validations/ updates
- Fetch requests to fetch state of the scene or particular objects, like:
    - fetch the game scene from the backend on initial boot
    - fetch components of an object on selection

> Note: Any state updates to the game scene should and must propagate through the backend