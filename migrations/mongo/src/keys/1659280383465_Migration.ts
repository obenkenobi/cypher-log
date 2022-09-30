import { Db } from 'mongodb'
import { MigrationInterface } from 'mongo-migrate-ts';

export class Migration1659280383465 implements MigrationInterface {
  public async up(db: Db): Promise<any> {
    await db.collection('users').createIndex({ userName: 1 },
        { unique: true, name: "idx-users-username-unique" })
  }

  public async down(db: Db): Promise<any> {
    await db.collection('users').dropIndex( "idx-users-username-unique" )
  }
}
