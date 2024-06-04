<!--forhugo
+++
title="MySQL Replication -  Configuration and Troubleshooting"
+++
forhugo-->
# MySQL Replication - Configuration and Troubleshooting

This project is designed to introduce you to MySQL database setup, schema creation, replication, and failover using Amazon EC2 instances.

You will set up a primary MySQL server, test it, add a secondary server for replication, and demonstrate replication and failover processes. Most steps are designed to succeed, but some will require some troubleshooting and problem-solving.

## Learning Objectives
- Install MySQL server
- Configure MySQL for replication
- Troubleshoot the replication configuration
- Fail-over the replicating primary

Timebox: 2 days

## Project

There are different ways to configure MySQL replication. In this exercise, you will be configuring your servers for primary-replica (or master-slave) replication.

In this type of replication, the primary server (or the master) takes all the writes and they are automatically replicated onto the replica server (or the slave). This technique is widely used to increase the scalability of the database for read-intensive operations (which is extremely common for the web). In the primary-replica setup, the primary would normally be used for writes and replica (or replicas) for reads only. Even though it's technically possible to use the primary for the reads and the writes, it is impossible to write directly to the replica.

Another advantage of using such a replication setup is database resilience. For example, it is recommended to setup primary-replica with one primary and two or three replicas in different availability zones. In the event of one availability zone (or datacenter) going down, the database will continue functioning flawlessly as other replicas will be used for reading. In a different scenario of one replica crashing, it can be replaced while the remaining replicas are serving the reads. Should the primary crash, an operation called a 'fail-over' should be carried out: one replica is promoted to be a primary while another MySQL server is being stood up in place of a broken primary.

### Task 1: Set Up the Primary MySQL Server

1. **Launch an EC2 Instance**
   - Name your instance `db-proj-<yourname>-primary`
   - Choose an Ubuntu AMI
   - Select a `t2.micro` instance type
   - Use key pair if you have one (recommended)
   - Select existing security group:
     - allow ssh access
     - open-database-port-internal
   - Keep default settings for storage
   - Review and launch the instance
   - Connect to the instance using SSH

2. **Install MySQL Server**
   - Update the package list and install MySQL server:
     ```bash
     sudo apt-get update
     sudo apt-get install mysql-server
     ```
   - Secure the MySQL installation:
     Answer `Y` for everything, and for this exercise choose *LOW* password validation policy)
     ```bash
     sudo mysql_secure_installation
     ```

3. **Configure MySQL for Remote Access**
   - Edit the MySQL configuration file to allow remote connections:
     ```bash
     sudo nano /etc/mysql/mysql.conf.d/mysqld.cnf
     ```
   - Change `bind-address` to `0.0.0.0`.
   - Restart MySQL service:
     ```bash
     sudo systemctl restart mysql
     ```

4. **Create a MySQL User for Replication**
   - Log in to MySQL and create a user for replication:
     ```bash
     sudo mysql -u root
     ```
   - Run the following SQL commands:
     ```sql
     CREATE USER 'replica_user'@'%' IDENTIFIED BY 'C90L6`!Doe{K';
     GRANT REPLICATION SLAVE ON *.* TO 'replica_user'@'%';
     FLUSH PRIVILEGES;
     ```

5. **Create a Sample Database Schema**
   - Still in the MySQL console, create a test database and table:
     ```sql
     CREATE DATABASE cyfdb;
     USE cyfdb;
     CREATE TABLE users (
         id INT AUTO_INCREMENT PRIMARY KEY,
         name VARCHAR(100),
         email VARCHAR(100)
     );
     ```

### Task 2: Set Up the Secondary MySQL Server

1. **Launch a Second EC2 Instance**
   - Follow similar steps as for the primary instance, but tag it as `db-proj-<yourname>-replica` just for your reference.

2. **Install MySQL Server on the Secondary Instance**
   - Connect to the secondary instance using SSH.
   - Repeat the installation steps as for the primary server.

3. **Configure MySQL for Replication on the Secondary Server**
   - Stop MySQL service:
     ```bash
     sudo systemctl stop mysql
     ```
   - Edit the MySQL configuration file:
     ```bash
     sudo nano /etc/mysql/mysql.conf.d/mysqld.cnf
     ```
   - Add the following lines at the end:
     ```ini
     server-id=2
     relay-log=/var/log/mysql/mysql-relay-bin.log
     log_bin=/var/log/mysql/mysql-bin.log
     ```
   - Start MySQL service:
     ```bash
     sudo systemctl start mysql
     ```

4. **Configure Replication**
   - Obtain the master status on the primary server:
     ```sql
     SHOW MASTER STATUS;
     ```
   - Note down the `File` and `Position` values, also note down the private instance IP address (e.g. from your AWS console)
   - On the secondary server, set up the replication:
     ```sql
     CHANGE MASTER TO
         MASTER_HOST='<primary_server_ip>',
         MASTER_USER='replica_user',
         MASTER_PASSWORD='C90L6`!Doe{K',
         MASTER_LOG_FILE='<file_name_from_master_status>',
         MASTER_LOG_POS=<position_from_master_status>;
     START SLAVE;
     ```

5. **Verify Replication**
   - Check the replica status:
     ```sql
     SHOW REPLICA STATUS\G;
     ```
   - Ensure `Slave_IO_Running` and `Slave_SQL_Running` are both `Yes`.
   - Something isn't quite right. Can you figure out how to fix this? There may be a few things that need to be fixed. Use `SHOW REPLICA STATUS\G` to find out what is wrong. Can you see the `cyfdb` database on the replica? Keep a log of all commands you are executing on each server while troubleshooting this.

### Task 3: Demonstrate Replication

1. **Insert Data into the Primary Server**
   - On the primary server, insert a new record:
     ```sql
     INSERT INTO cyfdb.users (name, email) VALUES ('John Doe', 'john@example.com');
     ```

2. **Verify Data on the Secondary Server**
   - On the secondary server, query the table:
     ```sql
     SELECT * FROM testdb.users;
     ```
   - Verify that the data matches the primary server.
   - Before you execute the next statement, please write down what result you would expect from it:
   ```sql
   INSERT INTO cyfdb.users (name, email) VALUES ('Jane Doe', 'jane@example.com');
   ```
   Does the result match your prediction? If not, can you guess why?

### Task 4: Demonstrate Failover

1. **Simulate Primary Server Failure**
We're going to stop the primary server, to simulate some real failure (e.g. hardware failure or loss of network).
   - Stop the MySQL service on the primary server:
     ```bash
     sudo systemctl stop mysql
     ```

2. **Promote Secondary Server to Primary**
   - On the secondary server, stop the replica:
     ```sql
     STOP REPLICA;
     ```
   - Reset the replica configuration:
     ```sql
     RESET REPLICA ALL;
     ```
   - Ensure the secondary server can accept writes by setting the read-only mode to off:
     ```sql
     SET GLOBAL read_only = OFF;
     ```

3. **Test Write Operations on the New Primary**
   - Insert new data into the secondary server (now acting as the primary):
     ```sql
     INSERT INTO cyfdb.users (name, email) VALUES ('Jane Doe', 'jane@example.com');
     ```
   - Query the table to ensure the new data is inserted:
     ```sql
     SELECT * FROM cyfd.users;
     ```
4. **Service Location**
   - Write down your thoughts on how primary and replica can be conveniently located by their clients (e.g. a web application), given the fact that primary may be failed over and replicas replaced at any moment, and the new instances will receive a different IP address.
     
### Task 5: Reconfigure Original Primary as Secondary

1. **Reconfigure the Original Primary**
   - Start the MySQL service on the original primary:
     ```bash
     sudo systemctl start mysql
     ```
   - On the original primary server, set it up as the slave to the new primary:
     ```sql
     CHANGE MASTER TO
         MASTER_HOST='<new_primary_ip>',
         MASTER_USER='replica_user',
         MASTER_PASSWORD='C90L6`!Doe{K',
         MASTER_LOG_FILE='<file_name_from_new_primary>',
         MASTER_LOG_POS=<position_from_new_primary>;
     START SLAVE;
     ```
   - Check the replica status:
     ```sql
     SHOW REPLICA STATUS\G;
     ```

2. **Verify Data Synchronization**
   - Insert new data on the new primary and verify it replicates to the original primary.

   This step may again not work straight away. Refer to how you configured the original primary: did you follow every step for this server?

### Task 6: Add another replica (Optional)

Can you add another replica to this MySQL cluster using the steps above? Describe the issues you are having; what do you think might be the solution?
