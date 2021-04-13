make

# Add a blank line between make and run
echo

# Give the search path for the library
# LD_LIBRARY_PATH is for Linux, DYLD_LIBRARY_PATH for Mac
LD_LIBRARY_PATH=../c_binding DYLD_LIBRARY_PATH=../c_binding ./main
