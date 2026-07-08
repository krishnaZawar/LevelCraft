# Game Editor

The game editor is the place where users can built games without writing a single line of code, via drag and drop mechanisms and component customizations.

The editor should consist of the following panels to be able to provide a seamless game building experience:
- The hierarchy panel
- The attributes panel
- The workspace panel
- The utility panel

# Understanding Panels

## Hierarchy Panel

The hierarchy panel is used mainly to glance at the gameObjects, and their children, present in the game scene.

It also helps in selection, deselection, addition, deletion of gameObjects.

**Primary tasks would include:**
- display game scene objects and their hierarchies
- select/deselect objects

**future extensions include:**
- addition and deletion of gameObjects
- grouping of gameObjects to form complex hierarchies

## Attributes Panel

This is the panel which will be used to modify/update a gameObjects abilites

It allows for addition, deletion, updation of components 
for individual gameObjects. It also allows for updation of object metadata like name, tag, etc.

**Primary tasks would include:**
- display all the components and metadata of the gameObject
- add, delete, edit gameObject components
- update gameObject metadata

> Note: The update should only affect a single gameObject

## Workspace panel

This is the place where you can visualize the game scene, drag and drop components according to your need, update the game scene as desired.

**Primary tasks would include:**
- display the game scene in its intended visualized form
- drag and drop gameObjects for positioning
- selection/deselection of gameObjects

**future extensions include:**
- addition and deletion of gameObjects
- scaling of gameObjects
- group selection features

## Utility panel

The utility panel is used to add, delete, duplicate gameObjects. We should also be able to build and run game scenes from this panel.

**Primary tasks would include:**
- Addition, deletion, duplication of gameObjects
- build, run, testing of game scenes