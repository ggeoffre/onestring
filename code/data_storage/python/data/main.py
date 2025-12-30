# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025 ggeoffre, LLC

from cassandra_data import raw_cassandra_data
from mongo_data import raw_mongo_data

print("\nCASSANDRA\n#########")
raw_cassandra_data()

print("\nMONGO\n#########")
raw_mongo_data()
