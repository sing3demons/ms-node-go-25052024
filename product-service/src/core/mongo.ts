import { MongoClient } from 'mongodb'

export async function connectMongo() {
  const client = new MongoClient('mongodb://mongo1:27017,mongo2:27018,mongo3:27019/?replicaSet=my-replica-set')
  await client.connect()
  console.log('Connected to MongoDB')

  return client
}
