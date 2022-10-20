import { Db } from 'mongodb'
import { MigrationInterface } from 'mongo-migrate-ts';

export class Migration1664696511965 implements MigrationInterface {// userKeys
  public async up(db: Db): Promise<any> {
    await db.collection('userKeys').createIndex({ userId: 1 },
        { unique: true, name: "idx-userKeys-userId-unique" })
  }

  public async down(db: Db): Promise<any> {
    await db.collection('userKeys').dropIndex( "idx-userKeys-userId-unique" )
  }
}
