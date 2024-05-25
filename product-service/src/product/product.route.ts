import { type Router } from 'express'
import { IRoute, MyRouter } from '../core/my-route'
import { LoggerType } from '../core/logger'
import ProductService from './product.service'
import ProductController from './product.controller'
import { type MongoClient } from 'mongodb'

export class ProductRouter {
  constructor(private readonly route: IRoute, private readonly logger: LoggerType, private readonly client: MongoClient) {}

  register(): Router {
    const productService = new ProductService(this.client)
    const productController = new ProductController(this.route, this.logger, productService)
    return new MyRouter().Register(productController).instance
  }
}
