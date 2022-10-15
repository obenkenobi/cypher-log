import { Db } from 'mongodb'
import { MigrationInterface } from 'mongo-migrate-ts';

export class Migration1659279640978 implements MigrationInterface {
  public async up(db: Db): Promise<any> {
    await db.collection('users').createIndex({ authId: 1 },
        { unique: true, name: "idx-users-authId-unique" })
  }

  public async down(db: Db): Promise<any> {
    await db.collection('users').dropIndex( "idx-users-authId-unique" )
  }
}
