CREATE TABLE public.app_app (
                                app_id bigserial NOT NULL,
                                app_name varchar(400) NOT NULL,
                                app_status smallint NOT NULL,
                                time_created bigint NOT NULL,
                                time_updated bigint NOT NULL,
                                CONSTRAINT app_app_pk PRIMARY KEY (app_id)
);
CREATE UNIQUE INDEX app_app_app_name_idx ON public.app_app (app_name);

CREATE TABLE public.app_token (
                                  token_id bigserial NOT NULL,
                                  app_id bigint NOT NULL,
                                  "token" varchar(400) NOT NULL,
                                  security_value varchar(400) NULL,
                                  token_type smallint NOT NULL,
                                  token_status smallint NOT NULL,
                                  expirydate_millis bigint NULL,
                                  time_created bigint NOT NULL,
                                  time_updated bigint NOT NULL,
                                  CONSTRAINT app_token_pk PRIMARY KEY (token_id)
);
CREATE UNIQUE INDEX app_token_token_idx ON public.app_token ("token");

INSERT INTO app_app(APP_NAME, APP_STATUS, TIME_CREATED, TIME_UPDATED) VALUES ('user', 0, 1, 1);

INSERT INTO app_token(app_id, token, token_type, token_status, expirydate_millis, time_created, time_updated)
    VALUES (1, 'userservtoken', 1, 0, 253373578095000, 1, 1);
