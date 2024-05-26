import axios, { AxiosResponse, AxiosRequestConfig, RawAxiosRequestHeaders } from 'axios'
import type { Request, Response, NextFunction } from 'express'
import jwt from 'jsonwebtoken'
import Config from '../config'

type IAuthService = {
  message: string
  success: string
  data: {
    access_token: string
  }
}

type IBody = {
  access_token: string
}

type PayloadToken = {
  sub: string
  exp: number
  iat: number
  email: string
}

type CustomRequest = Request & { token?: PayloadToken }

export default class AuthService {
  validateToken = async (req: Request, res: Response, next: NextFunction) => {
    const authorization = req.headers['authorization']
    if (!authorization) {
      res.status(401).send('Unauthorized')
      return
    }
    const access_token = authorization.split(' ')[1]

    const config: AxiosRequestConfig = { headers: { 'Content-Type': 'application/json' } }
    const data = { access_token: access_token }

    try {
      const url = Config.AUTH_SERVICE
      const result = await axios.post<IBody, AxiosResponse<IAuthService>, IBody>(url, data, config)

      if (result.status === 200) {
        const token = result.data.data.access_token

        const payload = jwt.decode(token) as PayloadToken

        ;(req as CustomRequest).token = payload

        next()
      } else {
        res.status(401).send('Unauthorized')
      }
    } catch (error) {
      res.status(401).send('Unauthorized')
    }
  }
}
