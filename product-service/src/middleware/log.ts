import type { Request, Response, NextFunction } from 'express'
import { v4 as uuidv4 } from 'uuid'
import { ILogger } from '../core/logger'

const httpLogger = (logger: ILogger) => (req: Request, res: Response, next: NextFunction) => {
  const start = performance.now()
  const session = req.header('x-session') ?? `default-${uuidv4()}`
  const data = {
    request: {
      headers: req.headers,
      host: req.headers.host,
      baseUrl: req.baseUrl,
      url: req.url,
      method: req.method,
      body: req.body,
      params: req?.params,
      query: req?.query,
      clientIp: req.headers['x-forwarded-for'] ?? req?.socket.remoteAddress,
    },
    response: {
      headers: res.getHeaders(),
      statusCode: res.statusCode,
      body: {} as any,
    },
  }

  const originalSend = res.json
  res.json = (body) => {
    data.response.body = body
    return originalSend.call(res, body)
  }
  next()
  res.on('finish', () => {
    const end = performance.now()
    const logData = { session, ...data, duration: (end - start).toFixed(2) }
    logger.info('summary', logData)
  })
}

export default httpLogger
