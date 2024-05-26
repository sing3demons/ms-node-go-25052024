import { z } from 'zod'
import Context from '../core/context'
import { LoggerType } from '../core/logger'
import { IRoute } from '../core/my-route'
import { IProduct, IProductQuerySchema, ProductBodySchema } from './product.model'
import ProductService from './product.service'
import config from '../config'

export default class ProductController {
  private readonly host = config.HOST
  constructor(
    private readonly route: IRoute,
    private readonly logger: LoggerType,
    private readonly productService: ProductService
  ) {}

  //   page?: string
  //   pageSize?: string
  //   sort?: string
  //   order?: string

  getProducts = this.route
    .get('/')
    .query(IProductQuerySchema)
    .handler(async ({ query }) => {
      const ctx = Context.get()
      const logger = this.logger.Logger(ctx)
      logger.info('Get products')
      const { data, total } = await this.productService.getProducts(logger, query)
      const result: IProduct[] = data.map((item) => {
        return {
          id: item?.id,
          href: item?.id && `${this.host}/products/${item?.id}`,
          name: item?.name,
          description: item?.description,
          language: item?.language,
          price: item?.price,
        }
      })
      return {
        data: result,
        total,
        page: query.page ? parseInt(query.page) : 1,
        pageSize: query.pageSize ? parseInt(query.pageSize) : 10,
      }
    })

  getProductsById = this.route
    .get('/:id')
    .params(z.object({ id: z.string() }))
    .handler(async ({ params }) => {
      const ctx = Context.get()
      const logger = this.logger.Logger(ctx)
      logger.info('Get product by id')
      const data = await this.productService.getProductById(logger, params.id)

      if (!data) {
        return {
          message: 'Product not found',
          success: false,
          statusCode: 404,
        }
      }
      data.href = `${this.host}/products/${data.id}`

      return {
        data,
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
