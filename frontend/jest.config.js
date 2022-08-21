module.exports = {
  roots: ['<rootDir>'],
  testMatch: [
    // '***/tests/**/*.js?(x)',
    '**/?(*.)(spec|test).(js|ts|tsx)',
  ],
  transform: {
    '^.+\\.(ts|tsx)$': 'ts-jest',
  },
};
