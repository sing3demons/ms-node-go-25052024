import { type Router } from 'express'
import ProductService from './product.service'
import ProductController from './product.controller'
import { type MongoClient } from 'mongodb'
import { IRoute, LoggerType, MyRouter } from '@express-zod/sing3demons'

export class ProductRouter {
  constructor(
    private readonly route: IRoute,
    private readonly logger: LoggerType,
    private readonly client: MongoClient
  ) {}

  register(): Router {
    const productService = new ProductService(this.client)
    const productController = new ProductController(this.route, this.logger, productService)
    return new MyRouter().Register(productController).instance
  }
}
