import { Router } from 'express'
import { ProductRouter } from './product/product.route'
import { type MongoClient } from 'mongodb'
import AuthService from './middleware/auth'
import { IRoute, LoggerType, TypeRoute } from '@express-zod/sing3demons'

export default function registerRoute(myRoute: IRoute = new TypeRoute(), client: MongoClient, logger: LoggerType) {
  const router = Router()
  const { validateToken } = new AuthService()

  router.get('/healthz', (_req, res) => res.status(200).json('OK'))

  router.use('/products', validateToken, new ProductRouter(myRoute, logger, client).register())
  return router
}
