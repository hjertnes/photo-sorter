# photo-sorter

Just a simple utility to take a folder with photos and rename them based on dates in their meta data (either exif or file system) with a fallback to "now" if neither is found

It has three modes:
- default is to take a dry run this writes nothing
- --flat renames into a flat structure with all the files at root in output dir
- --nested nests the files in folders for year and month and then the same filename as above.

## File name structures:
- Flat: yyyy-mm-dd hh:mm:ss originalfilename
- Nested: yyyy/mm/yyyy-mm-dd hh:mm:ss originalfilename