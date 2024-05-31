import express from 'express'
import registerRoute from './root-route'
import { globalErrorHandler, TypeRoute, Context, Logger } from '@express-zod/sing3demons'
import { connectMongo } from './core/mongo'
import config from './config'
import httpLogger from './middleware/log'

async function main() {
  const myRoute = new TypeRoute()
  const logger = new Logger()

  const db = await connectMongo()
  const app = express()
  app.use(Context.Ctx)
  app.use(httpLogger(logger))
  app.use(express.json())
  app.use(express.urlencoded({ extended: true }))

  app.use('/', registerRoute(myRoute, db, logger))

  app.use(globalErrorHandler)

  app.listen(config.PORT, () => {
    logger.info('Server is running on port ' + config.PORT)
  })
}

main().catch((err) => {
  console.error(err)
})
