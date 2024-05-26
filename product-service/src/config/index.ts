import * as dotenv from 'dotenv'
const NODE_ENV = process.env.NODE_ENV || 'development'

if (NODE_ENV !== 'production') {
  dotenv.config({ path: '.env.dev' })
}

class Config {
  public readonly PORT: string
  public readonly HOST: string
  public readonly MONGO_URL: string
  public readonly AUTH_SERVICE: string

  constructor() {
    this.PORT = process.env.PORT || '3001'
    this.HOST = process.env.HOST || 'http://localhost:' + this.PORT
    this.MONGO_URL = process.env.MONGO_URL!
    this.AUTH_SERVICE = process.env.AUTH_SERVICE!

    console.log('loaded config ...')
  }
}

export default new Config()
