import { z } from 'zod'
import Context from '../core/context'
import { LoggerType } from '../core/logger'
import { IRoute } from '../core/my-route'
import { IProduct, ProductBodySchema } from './product.model'
import ProductService from './product.service'

export default class ProductController {
  constructor(
    private readonly route: IRoute,
    private readonly logger: LoggerType,
    private readonly productService: ProductService
  ) {}

  getProducts = this.route.get('/').handler(async ({}) => {
    const ctx = Context.get()
    const logger = this.logger.Logger(ctx)
    logger.info('Get products')
    const data = await this.productService.getProducts(logger)
    const result: IProduct[] = data.map((item) => {
      return {
        id: item?.id,
        href: item?.id && `/products/${item?.id}`,
        name: item?.name,
        description: item?.description,
        language: item?.language,
        price: item?.price,
      }
    })
    return {
      data: result,
    }
  })

  create = this.route
    .post('/')
    .body(ProductBodySchema)
    .handler(async ({ body }) => {
      const ctx = Context.get()
      const logger = this.logger.Logger(ctx)
      logger.info('Create product')

      const data = await this.productService.createProduct(logger, body)

      return {
        message: !data ? 'Create product failed' : 'Create product success',
        success: !!data,
        statusCode: !data ? 400 : 201,
        data: data,
      }
    })
}
