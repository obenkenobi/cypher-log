module.exports = {
  async up(db, client) {
    await db.collection('users').createIndex({ userName: 1 },
        { unique: true, name: "idx-users-username-unique" })
    await db.collection('users').createIndex({ authId: 1 },
        { unique: true, name: "idx-users-authId-unique" })
  },

  async down(db, client) {
    await db.collection('users').dropIndexes( ["idx-users-username-unique", "idx-users-authId-unique"] )
  }
};
