# SR Spray Viewer Project

## Ussage
 - Download the precompiled zip folder and drop it in your Garry's Mod data folder which should look something like this `C:\Program Files (x86)\Steam\steamapps\common\GarrysMod\garrysmod\data`
 - Run 'build.exe' it should take a second to do it's thing and hopefully return no errors
 - Open viewsprays.html, if it opens by a text editor by deafult use open with in explorer and select your perfered browser
 ## Progress
  - I got a lot of the main features of the build.go file done as of now 2023-02-09 it's able to convert the spray folder into jpg(s) succesfully and it's also able to create a table of each steamid and their sprays
  - things that still need to be worked on are as following
   - I still need to find a way to parse the cached avatars at `\webmaterial\avatars\` so that I can use them for the frontend
   - Front end still defiently needs to be worked on, as of now it doesnt work and it also needs to use `https://github.com/12pt/steamid-converter` to convert steamID64s into steamIDs so that I can use `\sr_sprays\sr_playerdb.json` for player names
