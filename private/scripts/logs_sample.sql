INSERT INTO logs (_timestamp,
                  _namespace,
                  host,
                  service,
                  level,
                  user_id,
                  session_id,
                  trace_id,
                  _source,
                  string_names,
                  string_values,
                  int_names,
                  int_values,
                  float_names,
                  float_values,
                  bool_names,
                  bool_values)
VALUES (now(),
        'prod',
        'host-01',
        'auth-service',
        'info',
        12345,
        'session-xyz',
        'trace-abc',
        '{"endpoint":"/login","user_id":12345,"success":true,"latency":42.7}',

           -- string dynamic fields
        ['endpoint'],
        ['/login'],

           -- int dynamic fields
        ['user_id'],
        [12345],

           -- float dynamic fields
        ['latency'],
        [42.7],

           -- boolean dynamic fields (stored as strings)
        ['success'],
        ['true']);
