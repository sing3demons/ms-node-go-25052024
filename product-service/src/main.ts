import express from 'express'
import registerRoute from './root-route'
import { globalErrorHandler, TypeRoute } from './core/my-route'
import Logger from './core/logger'
import Context from './core/context'
import { connectMongo } from './core/mongo'

async function main() {
  const myRoute = new TypeRoute()
  const logger = new Logger()

  const db = await connectMongo()
  const app = express()
  app.use(Context.Ctx)
  app.use(express.json())
  app.use(express.urlencoded({ extended: true }))

  app.use('/', registerRoute(myRoute, db, logger))

  app.use(globalErrorHandler)

  app.listen(3000, () => {
    logger.info('Server is running on port 3000')
  })
}

main().catch((err) => {
  console.error(err)
})
