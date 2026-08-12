-- @table DataOrigin
DataOriginID iNTEGER PRIMARY KEY
    CHECK (DataOriginID >= 0),
Name TEXT NOT NULL,
Description TEXT 
-- DO NOT ADD METADATA AT THE END 

