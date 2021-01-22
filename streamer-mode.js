// streamer-mode developed for the most EXTREME alpha member streams
// in the hopes of not leaking anything when using OBS

const DataFileEx = require("datafile_ex.js");

Cheat.ExecuteCommand("hideconsole");

// We don't store the data directly into key to be able to parse data easier
var dataCache = {
    "CONSOLE": false,
    "MENU": {
        "OPEN": UI.IsMenuOpen(), // BOOL
        "INFO": UI.GetMenuPosition(), // [X, Y, W, H]
    }
};

var dataFile = "streamermode.data";

var CM = false;

function CreateMove() {
    CM = true;
    Run();
}

function Draw() {
    if (!CM) {
        Run();
    }
}

function Run() {
    // 192 is the ~ key (csgo console)
    /*
        if (Input.IsKeyPressed(192) && !dataCache["CONSOLE"]) {
        Cheat.Print("test")
            dataCache["CONSOLE"] = true;
            DataFileEx.save(dataFile, dataCache)
        } else if (Input.IsKeyPressed(192) && dataCache["CONSOLE"]) {
            dataCache["CONSOLE"] = false;
            DataFileEx.save(dataFile, dataCache);
        }
    */
    if (UI.IsMenuOpen() && !dataCache["MENU"]["OPEN"]) {
        dataCache["MENU"]["OPEN"] = true;
        DataFileEx.save(dataFile, dataCache)
    } else if (!UI.IsMenuOpen() && dataCache["MENU"]["OPEN"]) {
        dataCache["MENU"]["OPEN"] = false;
        DataFileEx.save(dataFile, dataCache);
    }

    if (dataCache["MENU"]["INFO"].toString() !== UI.GetMenuPosition().toString() && UI.IsMenuOpen()) {
        dataCache["MENU"]["INFO"] = UI.GetMenuPosition();
        DataFileEx.save(dataFile, dataCache);
    }
}

function Unload() {
    DataFileEx.save(dataFile, dataCache);
}

Cheat.RegisterCallback("CreateMove", CreateMove.name);
Cheat.RegisterCallback("Draw", Draw.name);
Cheat.RegisterCallback("Unload", Unload.name);