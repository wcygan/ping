# Introduction

In this application, we want to store each individual ping that is sent to our API. We will store the pings in a Postgres database.

This guide aims to provide a step-by-step guide to setting up a Postgres database on a Kubernetes cluster using CNPG and Helm, with the installation and deployment orchestrated using Skaffold.

## Table of Contents

1. Install Helm
2. Add CNPG Helm Chart Repository
3. Install the CNPG Operator
4. Deploy a CNPG Cluster
5. Connect to Postgres
6. Create a Ping Table
7. Implement the API logic to store and retrieve pings

## 1: Install Helm and hook it up using Skaffold

### Install Helm

```bash
brew install helm
```

## 2: Add CNPG Helm Chart Repository

## 3: Install the CNPG Operator

### 3.1: Configure Skaffold to use Helm

## 4: Deploy a CNPG Cluster

## 5: Connect to Postgres

## 6: Create a Ping Table

We only want to store the ID of the ping and the timestamp of when the ping was initiated (provided by the user's system clock).

We will use a single table to store the pings.

The table will have the following columns:

- id (UUID)
- pinged_at (timestamp)

The SQL for the table is as follows:

```sql
CREATE TABLE pings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pinged_at TIMESTAMP WITH TIME ZONE
);
```

The SQL for inserting a ping is as follows:

```sql
INSERT INTO pings (id, pinged_at) VALUES (gen_random_uuid(), '2023-10-05 14:30:00+00');
```

## 7: Implement the API logic to store and retrieve pings