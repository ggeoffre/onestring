# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

from cassandra_data import raw_cassandra_data
from mongo_data import raw_mongo_data
from mysql_data import raw_mysql_data
from postgres_data import raw_postgres_data
from redis_data import raw_redis_data

print("\nCASSANDRA\n#########")
raw_cassandra_data()

print("\nMONGO\n#########")
raw_mongo_data()

print("\nMYSQL\n#########")
raw_mysql_data()

print("\nPOSTGRES\n#########")
raw_postgres_data()

print("\nREDIS\n#########")
raw_redis_data()
