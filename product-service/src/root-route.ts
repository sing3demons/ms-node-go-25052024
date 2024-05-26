import { Router } from 'express'
import { IRoute, TypeRoute } from './core/my-route'
import { ProductRouter } from './product/product.route'
import { LoggerType } from './core/logger'
import { type MongoClient } from 'mongodb'
import AuthService from './middleware/auth'

export default function registerRoute(myRoute: IRoute = new TypeRoute(), client: MongoClient, logger: LoggerType) {
  const router = Router()
  const { validateToken } = new AuthService()

  router.get('/healthz', (_req, res) => res.status(200).json('OK'))

  router.use('/products', validateToken, new ProductRouter(myRoute, logger, client).register())
  return router
}
