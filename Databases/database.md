# Databases

> A database is a systematic collection of data. 
> They support electronic storage and manipulation of data. Databases make data management easy.

Database main function:
- store data
- update and delete data
- return data according to a query
- administer the database

A database needs to be:
- reliable
- efficient
- correct

## CAP theorem

> The CAP theorem says that any distributed database can only satisfy two of the three features.

- Consistency: every node responds with the most recent version of the data
- Availability: any mode can send a response
- Partition Tolerance: the system continues working even if communication between any of the nodes is broken

Tradeoffs: 
Stability of the network is not guarantee, meaning a db should always satisfy Partition tolerance. 
BUT it implies a tradeoff between Consistency and Availability: 
1. CP: Satisfy concurrency and rollback unfinished operations and wait until all nodes are back
2. AP: Continue responding but risk inconsistencies

![img.png](01-cp-ap-databases.png)

## Transactions: ACID

> A transaction: series of database operations = single unit of work. They wither all succeed or all fail.

- Atomicity: all operations succeed or fail
- Consistency: a successful transaction leave the db valid with no schema violations
- Isolation: transactions can be executed concurrently
- Durability: A committed transaction is persisted in memory

Reality:
- Theoretically impossible to implement together
- Usually relational databases do support ACID transactions, and non-relational databases don’t

## Types of Databases
Here are some popular types of databases.

### Distributed databases
A distributed database is a type of database that has contributions from the common database and information captured by local computers. 
In this type of database system, the data is not in one place and is distributed at various organizations.

### Relational databases:
This type of database defines database relationships in the form of tables. 
It is also called Relational DBMS, which is the most popular DBMS type in the market. Database example of the RDBMS system include MySQL, Oracle, and Microsoft SQL Server database.

### Object-oriented relational databases:
This type of computers database supports the storage of all data types. The data is stored in the form of objects. 
The objects to be held in the database have attributes and methods that define what to do with the data. PostgreSQL is an example of an object-oriented relational DBMS.

### Centralized database:
It is a centralized location, and users from different backgrounds can access this data. 
This type of computers databases store application procedures that help users access the data even from a remote location.

### NoSQL databases:
NoSQL database is used for large sets of distributed data. There are a few big data performance problems that are effectively handled by relational databases. This type of computers database is very efficient in analyzing large-size unstructured data.
Graph databases:
A graph-oriented database uses graph theory to store, map, and query relationships. These kinds of computers databases are mostly used for analyzing interconnections. For example, an organization can use a graph database to mine data about customers from social media.

### OLTP(Online transaction processing) databases:
OLTP another database type which able to perform fast query processing and maintaining data integrity in multi-access environments.

> OLTP and OLAP: The two terms look similar but refer to different kinds of systems. Online transaction processing (OLTP) captures, stores, and processes data from transactions in real time.
> Online analytical processing (OLAP) uses complex queries to analyze aggregated historical data from OLTP systems.

### Document/JSON database:
In a document-oriented database, the data is kept in document collections, usually using the XML, JSON, BSON formats.
One record can store as much data as you want, in any data type (or types) you prefer.

