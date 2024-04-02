#gsb

Short for "Golang Service Bus", it provides a messaging system based on the Service Bus Pattern

## Messaging

Default is beanstalkd.
make run_beanstalkd

## SonarQube

SonarQube is used for static code analysis. Entries for running both the server and the cli are in the Makefile.

## Principles
1. Errors should give context information and be as specific as possible
2. Convention over configuration
3. If there are errors during start up configuration, fail
4. If there are errors after start up configuration, don't fail
5. User space should raise errors
6. Transactions are managed in GSB
