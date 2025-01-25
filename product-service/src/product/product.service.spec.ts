import { MongoClient } from 'mongodb'
import { ILogger } from '../core/logger'
import ProductService from './product.service'
import { mock } from 'jest-mock-extended'
import { IProductBody } from './product.model'
import { v4 as uuid } from 'uuid'

describe('ProductService', () => {
  let client = jest.mocked(MongoClient) as unknown as MongoClient
  const productService = new ProductService(client)
  const logger = mock<ILogger>()

  beforeEach(() => {
    const mockSession = jest.fn().mockReturnValue({
      startTransaction: jest.fn(),
      commitTransaction: jest.fn(),
      abortTransaction: jest.fn(),
      endSession: jest.fn(),
    })

    client.startSession = jest.fn().mockImplementation(mockSession)
    client.db = jest.fn().mockReturnValue({
      collection: jest.fn().mockReturnValue({
        insertOne: jest.fn().mockReturnValue({
          insertedId: '51226c6b-9857-47a5-be79-733685486a2d',
        }),
        insertMany: jest.fn(),
        find: jest.fn().mockReturnValue({
          toArray: jest.fn().mockReturnValue([]),
        }),
        countDocuments: jest.fn().mockReturnValue(0),
        findOne: jest.fn().mockReturnValue(null),
      }),
    })
  })

  afterEach(() => {
    jest.clearAllMocks()
  })

  it('createProduct', async () => {
    const body: IProductBody = {
      name: 'Product 1',
      price: 100,
    }
    const result = await productService.createProduct(logger, body)

    expect(client.startSession).toHaveBeenCalled()
    expect(client.db).toHaveBeenCalled()

    expect(result?.insertedId).toBe('51226c6b-9857-47a5-be79-733685486a2d')
  })

  it('createProduct 2', async () => {
    const body: IProductBody = {
      name: 'Product 1',
      price: 100,
      th: {
        name: 'Product 1',
        description: 'Description',
      },
    }
    const result = await productService.createProduct(logger, body)

    expect(client.startSession).toHaveBeenCalled()
    expect(client.db).toHaveBeenCalled()

    expect(result?.insertedId).toBe('51226c6b-9857-47a5-be79-733685486a2d')
  })

  it('createProduct - error', async () => {
    const body: IProductBody = {
      name: 'Product 1',
      price: 100,
    }
    client.db = jest.fn().mockImplementation(() => {
      throw new Error('Error')
    })

    await expect(productService.createProduct(logger, body)).rejects.toThrow('Error')
  })

  it('createProduct - error 2', async () => {
    const endSessionMock = jest.fn();
  const abortTransactionMock = jest.fn();
  const mockSession = {
    startTransaction: jest.fn(),
    commitTransaction: jest.fn(),
    abortTransaction: abortTransactionMock,
    endSession: endSessionMock,
  };

  // Mock `startSession`
  client.startSession = jest.fn().mockReturnValue(mockSession);

  const mockError = new Error('Database error');

  // Mock `db` to return collections
  client.db = jest.fn().mockReturnValue({
    collection: jest.fn().mockReturnValue({
      insertOne: jest.fn().mockRejectedValueOnce(mockError),
      insertMany: jest.fn(),
    }),
  });

  const body: IProductBody = {
    name: 'Product 1',
    price: 100,
  };

  // Use `await expect(...).rejects` to catch and test the error
  await expect(productService.createProduct(logger, body)).rejects.toThrow(mockError);

  expect(client.startSession).toHaveBeenCalled();
  expect(client.db).toHaveBeenCalledWith('product');
  expect(mockSession.abortTransaction).toHaveBeenCalled();
  expect(mockSession.endSession).toHaveBeenCalled();
  })

  it('getProducts - success return {data: [], total: 0}', async () => {
    const query = {
      pageSize: '10',
      page: '1',
      sort: 'name',
      order: 'asc',
    }
    const result = await productService.getProducts(logger, query)

    expect(client.db).toHaveBeenCalled()
    expect(result).toEqual({
      data: [],
      total: 0,
    })
  })

  it('getProductById - success', async () => {
    const id = uuid()
    const result = await productService.getProductById(logger, id)

    expect(client.db).toHaveBeenCalled()
    expect(result).toBeNull()
  })
})
