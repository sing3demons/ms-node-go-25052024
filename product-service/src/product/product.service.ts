import { ILogger } from '../core/logger'
import type { Collection, MongoClient, Filter, FindOptions } from 'mongodb'
import {
  IProduct,
  IProductBody,
  IProductQuery,
  IProductSchema,
  ProductLanguageSchema,
  ProductPriceLanguageSchema,
  ProductPriceSchema,
  ProductSchema,
} from './product.model'
import { z } from 'zod'
import { v4 as uuid } from 'uuid'

export default class ProductService {
  constructor(private readonly client: MongoClient) {}

  private getCollection<T extends object>(dbName: string, name: string): Collection<T> {
    return this.client.db(dbName).collection<T>(name)
  }

  async getProducts(logger: ILogger, query: IProductQuery) {
    const col = this.getCollection<IProductSchema>('product', 'product')
    logger.info(ProductService.name, {
      service: 'product-service',
      method: 'getProducts',
      message: 'Get products2',
    })

    const pageSize = query.pageSize ? parseInt(query.pageSize) : 10
    const page = query.page ? parseInt(query.page) : 1

    const filter: Filter<IProductSchema> = {}
    const options: FindOptions = {
      limit: pageSize,
      skip: (page - 1) * pageSize,
      sort: query.sort ? { [query.sort]: query.order === 'asc' ? 1 : -1 } : undefined,
    }

    const [data, total] = await Promise.all([col.find(filter, options).toArray(), col.countDocuments(filter)])
    return {
      data,
      total,
    }
  }

  async getProductById(logger: ILogger, id: string) {
    const col = this.getCollection<IProduct>('product', 'product')
    logger.info(ProductService.name, {
      service: 'product-service',
      method: 'getProductById',
      message: 'Get product by id',
    })

    return await col.findOne({ id })
  }

  async createProduct(logger: ILogger, body: IProductBody) {
    logger.info(ProductService.name, {
      service: 'product-service',
      method: 'createProduct',
      message: 'Create product',
    })

    const session = this.client.startSession()

    const productCol = this.getCollection('product', 'product')
    const productPriceCol = this.getCollection('product', 'price')
    const productLangCol = this.getCollection('productLanguage', 'productLanguage')
    const productLangPriceCol = this.getCollection('productLanguage', 'priceLanguage')

    try {
      session.startTransaction()
      const { name, description, price, th, en, vat, tax } = body

      const productBody: z.infer<typeof ProductSchema> = {
        id: uuid(),
        name,
        description,
      }

      if (price) {
        const priceBody: z.infer<typeof ProductPriceSchema> = {
          id: uuid(),
          value: price,
          language: [],
          tax: {
            unit: tax ? '฿' : undefined,
            value: tax && tax,
          },
          vat: {
            unit: vat ? '฿' : undefined,
            value: vat && vat,
          },
          unit: price ? '฿' : undefined,
        }

        productBody.price = {
          id: priceBody.id,
          value: priceBody.value,
        }

        const productPriceLangBody: z.infer<typeof ProductPriceLanguageSchema>[] = []
        if (th) {
          productPriceLangBody.push({
            id: uuid(),
            languageCode: 'th',
            name: th.name,
            description: th.description,
            unit: th.unit && th.unit,
            price: th.price,
            tax: {
              unit: tax ? '฿' : undefined,
              value: tax && tax,
            },
            vat: {
              unit: vat ? '฿' : undefined,
              value: vat && vat,
            },
          })

          const thPrice = productPriceLangBody.find((lang) => {
            if (lang.languageCode === 'th') {
              return {
                id: lang.id,
                languageCode: lang.languageCode,
                name: lang.name,
              }
            }
          })
          if (thPrice) {
            priceBody?.language?.push(thPrice)
          }
        }

        productPriceLangBody.push({
          id: uuid(),
          languageCode: 'en',
          name: en?.name ?? name,
          description: en?.description ?? description,
          unit: (en?.unit && en?.unit) ?? (price ? '฿' : undefined),
          price: en?.price ?? price,
          tax: {
            unit: tax ? '฿' : undefined,
            value: tax && tax,
          },
          vat: {
            unit: vat ? '฿' : undefined,
            value: vat && vat,
          },
        })

        const enPrice = productPriceLangBody.find((lang) => {
          if (lang.languageCode === 'en') {
            return {
              id: lang.id,
              languageCode: lang.languageCode,
              name: lang.name,
            }
          }
        })

        if (enPrice) {
          priceBody?.language?.push(enPrice)
        }

        if (productPriceLangBody.length !== 0) {
          await productLangPriceCol.insertMany(productPriceLangBody, { session })
        }

        await productPriceCol.insertOne(priceBody, { session })
      }

      const productLangBody: z.infer<typeof ProductLanguageSchema>[] = [
        {
          id: uuid(),
          languageCode: 'en',
          name: en?.name ?? name,
          description: en?.description ?? description,
        },
      ]

      if (th) {
        productLangBody.push({
          id: uuid(),
          languageCode: 'th',
          name: th.name,
          description: th.description,
        })
      }

      if (productLangBody.length !== 0) {
        await productLangCol.insertMany(productLangBody, { session })

        productBody.language = productLangBody.map((lang) => ({
          id: lang.id,
          languageCode: lang.languageCode,
          name: lang.name,
          description: lang.description,
        }))
      }

      const result = await productCol.insertOne(productBody, { session })

      await session.commitTransaction()
      return result
    } catch (error) {
      await session.abortTransaction()

      logger.error(ProductService.name, {
        message: 'Create product failed',
        error,
      })
      throw error; 
    } finally {
      await session.endSession()
    }
  }
}
