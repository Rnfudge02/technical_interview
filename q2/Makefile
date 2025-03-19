#Test harness for secure shredder implementation
#This is not comprehensive
#Usage: make [target]
#Default: all (run full test suite)
#Exit code reference:
#1: No arguments
#2: Stat error
#3: Directory
#4: Empty file
#5: Permission denied
#6: Other error opening
#7: User aborted
#8: Random data error
#9: Seek error
#10: Write error
#11: Sync error
#12: Remove error

.PHONY: all argument_test normal_tests nonexistent_test empty_test negative_test ambigious_test permission_test restore clean

#Harness settings
SHREDDER := /usr/local/go/bin/go run shredder.go
SAMPLES_DIR := ./samples
TEST_FILES := shoulddelete.txt shouldntdelete.txt empty.txt

#@ is used to supress command printing

#Main test sequence
all: argument_test restore normal_tests nonexistent_test directory_test empty_test negative_test ambigious_test permission_test
	@echo "All tests completed"

#Test missing arguments
argument_test:
	@echo "--- Testing argument handling ---"
	@$(SHREDDER); \
	if [ $$? -eq 1 ]; then \
		echo "Test 1 (No arguments) passed"; \
	else \
		echo "Test 1 (No arguments) failed"; \
	fi

#Normal operation tests
normal_tests: restore
	@echo "--- Testing normal operation ---"

#This is all one command, runs program with desired file and checks output.
	@yes | $(SHREDDER) $(SAMPLES_DIR)/shoulddelete.txt; \
	if [ $$? -eq 0 ]; then \
		echo "Test 1 (Normal) passed"; \
	else \
		echo "Test 1 (Normal) failed"; \
	fi

#Command for running and testing with invalid permissions.
	@$(SHREDDER) $(SAMPLES_DIR)/shouldntdelete.txt; \
	if [ $$? -eq 5 ]; then \
		echo "Test 2 (Permissions check) passed"; \
	else \
		echo "Test 2 (Permissions check) failed"; \
	fi

directory_test: restore
	@echo "--- Testing directory handling ---"

#Check behavior when passed directory
	@$(SHREDDER) $(SAMPLES_DIR)/directory; \
	if [ $$? -eq 3 ]; then \
		echo "Test 3 (Directory check) passed"; \
	else \
		echo "Test 3 (Directory check) failed"; \
	fi

#Nonexistent file test
nonexistent_test:
	@echo "--- Testing nonexistent file handling ---"

#Command for running and testing invalid file case
	@$(SHREDDER) ./doesnt_exist.txt; \
	if [ $$? -eq 2 ]; then \
		echo "Test 4 (Nonexistent file) passed"; \
	else \
		echo "Test 4 (Nonexistent file) failed"; \
	fi

#Empty file test
empty_test: restore
	@echo "--- Testing empty file handling ---"

#Command for running and testing with an empty file
	@$(SHREDDER) $(SAMPLES_DIR)/empty.txt; \
	if [ $$? -eq 4 ]; then \
		echo "Test 5 (Empty file) passed"; \
	else \
		echo "Test 5 (Empty file) failed"; \
	fi

#Negative user response test
negative_test: restore
	@echo "--- Testing negative user response ---"

#Command for running a test with negative user input
	@yes n | $(SHREDDER) $(SAMPLES_DIR)/shoulddelete.txt; \
	if [ $$? -eq 7 ]; then \
		echo "Test 6 (Negative response) passed"; \
	else \
		echo "Test 6 (Negative response) failed"; \
	fi

#Empty file test
ambigious_test: restore
	@echo "--- Testing abnormal user response ---"

#Command for running and testing with an empty file
	@echo "q" | $(SHREDDER) $(SAMPLES_DIR)/shoulddelete.txt; \
	if [ $$? -eq 7 ]; then \
		echo "Test 7 (Ambigious response) passed"; \
	else \
		echo "Test 7 (Ambigious response) failed"; \
	fi

#Permission test (run with sudo)
permission_test: restore
	@echo "--- Testing privileged access ---"

#Command for running and testing with appropriate privileges
	@yes | sudo $(SHREDDER) $(SAMPLES_DIR)/shouldntdelete.txt; \
	if [ $$? -eq 0 ]; then \
		echo "Test 8 (Privileged delete) passed"; \
	else \
		echo "Test 8 (Privileged delete) failed"; \
	fi

#Restore test files from backups
restore:
	@echo "Restoring test files..."
	@mkdir -p ./samples/directory
	@cp -p $(SAMPLES_DIR)/shoulddelete.txt.bak $(SAMPLES_DIR)/shoulddelete.txt
	@sudo cp -p $(SAMPLES_DIR)/shouldntdelete.txt.bak $(SAMPLES_DIR)/shouldntdelete.txt
	@touch $(SAMPLES_DIR)/empty.txt
	@sudo chmod 644 $(SAMPLES_DIR)/shouldntdelete.txt

#Clean up test artifacts
clean:
	@rm -f $(SAMPLES_DIR)/shoulddelete.txt $(SAMPLES_DIR)/empty.txt
	@sudo rm -f $(SAMPLES_DIR)/shouldntdelete.txt