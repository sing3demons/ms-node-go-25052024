import { MongoClient } from 'mongodb'
import config from '../config'

export async function connectMongo() {
  const client = new MongoClient(config.MONGO_URL)
  await client.connect()
  console.log('Connected to MongoDB')

  return client
}
