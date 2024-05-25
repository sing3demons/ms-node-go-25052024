import { Router } from 'express'
import { IRoute, TypeRoute } from './core/my-route'
import { ProductRouter } from './product/product.route'
import { LoggerType } from './core/logger'
import { type MongoClient } from 'mongodb'

export default function registerRoute(myRoute: IRoute = new TypeRoute(), client: MongoClient, logger: LoggerType) {
  const router = Router()

  router.get('/healthz', (_req, res) => res.status(200).json('OK'))

  router.use('/products', new ProductRouter(myRoute, logger, client).register())
  return router
}
