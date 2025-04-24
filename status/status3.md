# Astropath | Status Report #3

## Recap
During the last three weeks, we were able to achieve a working prototype for the startrail functionality on the commandline. In the last week, we took this a step further and got the progress to be visualy displayed on the GUI of the application. At the same time, other team members are working on designing the UI and coming up with a clean and effecient UX for the application aswell as a system which we are calling the "doctor-process" which will check the user's computer for required dependecies. 

## Tasks Completed

- Met with our project advisor **Ruben**. (2x)
- Functioning Lighten/Darken workflow for both **JPEG/TIFF** files. Currently does not have a UI to drive it or customize the features. But given a list of filepaths, will properly produce the correct image output for that workflow configuration. Implimented systems to notify frontend of progress along the way. **(Tyler)**
- Using figma, created a mockup of the Main-Menu, Upload, Settings and Info Screen. Started prototyping the workflow screen. **(Greg)**
- Implimented (seperate branch) the UI design changes for the main-menu route in code. Also created application logo for landing page. **(Greg)**
- Started research on **doctor-report** feature. Planning on using open-source tool called **dcraw** to convert .arw files to .tiff **(Nolan)**

For metrics:
- Our group has met collectively 13 times. 
- Met w/ Ruben twice.
- 12 'quality' commits on github.
- Functioning lighten/darken feature on cli w/jpeg/tiff support

## Success

Our biggest success is the completion of two full pages with a streamlined UI, along with the functionality outlined above, which validates the effectiveness of the lighten/darken feature. The main-menu looks good. And if you your app looks good, it is good. Ask FTX.

## Roadblocks
Research for ARW conversion has been a bit of a challenge. However, this is something which is should be completed soon. We had to refactor the processing code a few times to make it work with the event system properly. However, it now let's the backend communicate with the frontend super easily.
To overcome the event issues, we did a refactor of the module system. The ARW conversion challenge is still left to be tackeled. Once this is finished, we will be able to move onto the upload and sequencing/workflow sections of the application.

All of the challenges we were facing were due to time commitment and external factors. During the next three weeks, this should be less of a problem.

## Changes & Deviations
We have made no major changes to our plan since the last update.

## Confidence
- Tyler - 5
-  Nolan - 5
-  Greg - 5
-  Average - 5


## Group Dynamics
We are working well as a group and there are nomajor issues. The biggest challenge is finding time that fits all of our schedules and time to meet with Reuben.