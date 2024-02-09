--
-- PostgreSQL database cluster dump
--

SET default_transaction_read_only = off;

SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;

--
-- Drop databases (except postgres and template1)
--

DROP DATABASE test;




--
-- Drop roles
--

DROP ROLE test;


--
-- Roles
--

CREATE ROLE test;
ALTER ROLE test WITH SUPERUSER INHERIT CREATEROLE CREATEDB LOGIN REPLICATION BYPASSRLS PASSWORD 'SCRAM-SHA-256$4096:yUmZ5IXYoy10TRKMORAp7w==$B/mAnZnODaP2RWv0lc7e9cMh+vaeRsWg4DfHl+p1FUE=:1aOFTrPwBvCCOPTvgBa70JPSuYxKDzRPQLc/1YxDnAk=';

--
-- User Configurations
--








--
-- Databases
--

--
-- Database "template1" dump
--

--
-- PostgreSQL database dump
--

-- Dumped from database version 16.1
-- Dumped by pg_dump version 16.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

UPDATE pg_catalog.pg_database SET datistemplate = false WHERE datname = 'template1';
DROP DATABASE template1;
--
-- Name: template1; Type: DATABASE; Schema: -; Owner: test
--

CREATE DATABASE template1 WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';


ALTER DATABASE template1 OWNER TO test;

\connect template1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: DATABASE template1; Type: COMMENT; Schema: -; Owner: test
--

COMMENT ON DATABASE template1 IS 'default template for new databases';


--
-- Name: template1; Type: DATABASE PROPERTIES; Schema: -; Owner: test
--

ALTER DATABASE template1 IS_TEMPLATE = true;


\connect template1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: DATABASE template1; Type: ACL; Schema: -; Owner: test
--

REVOKE CONNECT,TEMPORARY ON DATABASE template1 FROM PUBLIC;
GRANT CONNECT ON DATABASE template1 TO PUBLIC;


--
-- PostgreSQL database dump complete
--

--
-- Database "postgres" dump
--

--
-- PostgreSQL database dump
--

-- Dumped from database version 16.1
-- Dumped by pg_dump version 16.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

DROP DATABASE postgres;
--
-- Name: postgres; Type: DATABASE; Schema: -; Owner: test
--

CREATE DATABASE postgres WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';


ALTER DATABASE postgres OWNER TO test;

\connect postgres

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: DATABASE postgres; Type: COMMENT; Schema: -; Owner: test
--

COMMENT ON DATABASE postgres IS 'default administrative connection database';


--
-- PostgreSQL database dump complete
--

--
-- Database "test" dump
--

--
-- PostgreSQL database dump
--

-- Dumped from database version 16.1
-- Dumped by pg_dump version 16.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: test; Type: DATABASE; Schema: -; Owner: test
--

CREATE DATABASE test WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';
CREATE DATABASE test_atlas WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';


ALTER DATABASE test OWNER TO test;

\connect test

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: continuity_containers; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.continuity_containers (
    id uuid NOT NULL,
    identity_id uuid,
    name character varying(255) NOT NULL,
    payload jsonb,
    expires_at timestamp without time zone NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    nid uuid
);


ALTER TABLE public.continuity_containers OWNER TO test;

--
-- Name: courier_message_dispatches; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.courier_message_dispatches (
    id uuid NOT NULL,
    message_id uuid NOT NULL,
    status character varying(7) NOT NULL,
    error json,
    nid uuid NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.courier_message_dispatches OWNER TO test;

--
-- Name: courier_messages; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.courier_messages (
    id uuid NOT NULL,
    type integer NOT NULL,
    status integer NOT NULL,
    body text NOT NULL,
    subject character varying(255) NOT NULL,
    recipient character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    template_type character varying(255) DEFAULT ''::character varying NOT NULL,
    template_data bytea,
    nid uuid,
    send_count integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.courier_messages OWNER TO test;

--
-- Name: identities; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.identities (
    id uuid NOT NULL,
    schema_id character varying(2048) NOT NULL,
    traits jsonb NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    nid uuid,
    state character varying(255) DEFAULT 'active'::character varying NOT NULL,
    state_changed_at timestamp without time zone,
    metadata_public jsonb,
    metadata_admin jsonb
);


ALTER TABLE public.identities OWNER TO test;

--
-- Name: identity_credential_identifiers; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.identity_credential_identifiers (
    id uuid NOT NULL,
    identifier character varying(255) NOT NULL,
    identity_credential_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    nid uuid,
    identity_credential_type_id uuid NOT NULL
);


ALTER TABLE public.identity_credential_identifiers OWNER TO test;

--
-- Name: identity_credential_types; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.identity_credential_types (
    id uuid NOT NULL,
    name character varying(32) NOT NULL
);


ALTER TABLE public.identity_credential_types OWNER TO test;

--
-- Name: identity_credentials; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.identity_credentials (
    id uuid NOT NULL,
    config jsonb NOT NULL,
    identity_credential_type_id uuid NOT NULL,
    identity_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    nid uuid,
    version integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.identity_credentials OWNER TO test;

--
-- Name: identity_recovery_addresses; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.identity_recovery_addresses (
    id uuid NOT NULL,
    via character varying(16) NOT NULL,
    value character varying(400) NOT NULL,
    identity_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    nid uuid
);


ALTER TABLE public.identity_recovery_addresses OWNER TO test;

--
-- Name: identity_recovery_codes; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.identity_recovery_codes (
    id uuid NOT NULL,
    code character varying(64) NOT NULL,
    used_at timestamp without time zone,
    identity_recovery_address_id uuid,
    code_type integer NOT NULL,
    expires_at timestamp without time zone DEFAULT '2000-01-01 00:00:00'::timestamp without time zone NOT NULL,
    issued_at timestamp without time zone DEFAULT '2000-01-01 00:00:00'::timestamp without time zone NOT NULL,
    selfservice_recovery_flow_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    nid uuid NOT NULL,
    identity_id uuid NOT NULL
);


ALTER TABLE public.identity_recovery_codes OWNER TO test;

--
-- Name: identity_recovery_tokens; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.identity_recovery_tokens (
    id uuid NOT NULL,
    token character varying(64) NOT NULL,
    used boolean DEFAULT false NOT NULL,
    used_at timestamp without time zone,
    identity_recovery_address_id uuid,
    selfservice_recovery_flow_id uuid,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    expires_at timestamp without time zone DEFAULT '2000-01-01 00:00:00'::timestamp without time zone NOT NULL,
    issued_at timestamp without time zone DEFAULT '2000-01-01 00:00:00'::timestamp without time zone NOT NULL,
    nid uuid,
    identity_id uuid NOT NULL,
    token_type integer DEFAULT 0 NOT NULL,
    CONSTRAINT identity_recovery_tokens_token_type_ck CHECK (((token_type = 1) OR (token_type = 2)))
);


ALTER TABLE public.identity_recovery_tokens OWNER TO test;

--
-- Name: identity_verifiable_addresses; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.identity_verifiable_addresses (
    id uuid NOT NULL,
    status character varying(16) NOT NULL,
    via character varying(16) NOT NULL,
    verified boolean NOT NULL,
    value character varying(400) NOT NULL,
    verified_at timestamp without time zone,
    identity_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    nid uuid
);


ALTER TABLE public.identity_verifiable_addresses OWNER TO test;

--
-- Name: identity_verification_codes; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.identity_verification_codes (
    id uuid NOT NULL,
    code_hmac character varying(64) NOT NULL,
    used_at timestamp without time zone,
    identity_verifiable_address_id uuid,
    expires_at timestamp without time zone DEFAULT '2000-01-01 00:00:00'::timestamp without time zone NOT NULL,
    issued_at timestamp without time zone DEFAULT '2000-01-01 00:00:00'::timestamp without time zone NOT NULL,
    selfservice_verification_flow_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    nid uuid NOT NULL
);


ALTER TABLE public.identity_verification_codes OWNER TO test;

--
-- Name: identity_verification_tokens; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.identity_verification_tokens (
    id uuid NOT NULL,
    token character varying(64) NOT NULL,
    used boolean DEFAULT false NOT NULL,
    used_at timestamp without time zone,
    expires_at timestamp without time zone NOT NULL,
    issued_at timestamp without time zone NOT NULL,
    identity_verifiable_address_id uuid NOT NULL,
    selfservice_verification_flow_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    nid uuid
);


ALTER TABLE public.identity_verification_tokens OWNER TO test;

--
-- Name: keto_relation_tuples; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.keto_relation_tuples (
    shard_id uuid NOT NULL,
    nid uuid NOT NULL,
    namespace character varying(200) NOT NULL,
    object uuid NOT NULL,
    relation character varying(64) NOT NULL,
    subject_id uuid,
    subject_set_namespace character varying(200),
    subject_set_object uuid,
    subject_set_relation character varying(64),
    commit_time timestamp without time zone NOT NULL,
    CONSTRAINT chk_keto_rt_uuid_subject_type CHECK ((((subject_id IS NULL) AND (subject_set_namespace IS NOT NULL) AND (subject_set_object IS NOT NULL) AND (subject_set_relation IS NOT NULL)) OR ((subject_id IS NOT NULL) AND (subject_set_namespace IS NULL) AND (subject_set_object IS NULL) AND (subject_set_relation IS NULL))))
);


ALTER TABLE public.keto_relation_tuples OWNER TO test;

--
-- Name: keto_uuid_mappings; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.keto_uuid_mappings (
    id uuid NOT NULL,
    string_representation text NOT NULL
);


ALTER TABLE public.keto_uuid_mappings OWNER TO test;

--
-- Name: networks; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.networks (
    id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.networks OWNER TO test;

--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.schema_migration (
    version character varying(48) NOT NULL,
    version_self integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO test;

--
-- Name: selfservice_errors; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.selfservice_errors (
    id uuid NOT NULL,
    errors jsonb NOT NULL,
    seen_at timestamp without time zone,
    was_seen boolean NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    csrf_token character varying(255) DEFAULT ''::character varying NOT NULL,
    nid uuid
);


ALTER TABLE public.selfservice_errors OWNER TO test;

--
-- Name: selfservice_login_flows; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.selfservice_login_flows (
    id uuid NOT NULL,
    request_url text NOT NULL,
    issued_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    active_method character varying(32) NOT NULL,
    csrf_token character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    forced boolean DEFAULT false NOT NULL,
    type character varying(16) DEFAULT 'browser'::character varying NOT NULL,
    ui jsonb,
    nid uuid,
    requested_aal character varying(4) DEFAULT 'aal1'::character varying NOT NULL,
    internal_context jsonb NOT NULL,
    oauth2_login_challenge uuid,
    oauth2_login_challenge_data text
);


ALTER TABLE public.selfservice_login_flows OWNER TO test;

--
-- Name: selfservice_recovery_flows; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.selfservice_recovery_flows (
    id uuid NOT NULL,
    request_url text NOT NULL,
    issued_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    active_method character varying(32),
    csrf_token character varying(255) NOT NULL,
    state character varying(32) NOT NULL,
    recovered_identity_id uuid,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    type character varying(16) DEFAULT 'browser'::character varying NOT NULL,
    ui jsonb,
    nid uuid,
    submit_count integer DEFAULT 0 NOT NULL,
    skip_csrf_check boolean DEFAULT false NOT NULL
);


ALTER TABLE public.selfservice_recovery_flows OWNER TO test;

--
-- Name: selfservice_registration_flows; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.selfservice_registration_flows (
    id uuid NOT NULL,
    request_url text NOT NULL,
    issued_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    active_method character varying(32) NOT NULL,
    csrf_token character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    type character varying(16) DEFAULT 'browser'::character varying NOT NULL,
    ui jsonb,
    nid uuid,
    internal_context jsonb NOT NULL,
    oauth2_login_challenge uuid,
    oauth2_login_challenge_data text
);


ALTER TABLE public.selfservice_registration_flows OWNER TO test;

--
-- Name: selfservice_settings_flows; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.selfservice_settings_flows (
    id uuid NOT NULL,
    request_url text NOT NULL,
    issued_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    identity_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    active_method character varying(32),
    state character varying(255) DEFAULT 'show_form'::character varying NOT NULL,
    type character varying(16) DEFAULT 'browser'::character varying NOT NULL,
    ui jsonb,
    nid uuid,
    internal_context jsonb NOT NULL
);


ALTER TABLE public.selfservice_settings_flows OWNER TO test;

--
-- Name: selfservice_verification_flows; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.selfservice_verification_flows (
    id uuid NOT NULL,
    request_url text NOT NULL,
    issued_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    csrf_token character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    type character varying(16) DEFAULT 'browser'::character varying NOT NULL,
    state character varying(255) DEFAULT 'show_form'::character varying NOT NULL,
    active_method character varying(32),
    ui jsonb,
    nid uuid,
    submit_count integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.selfservice_verification_flows OWNER TO test;

--
-- Name: session_devices; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.session_devices (
    id uuid NOT NULL,
    ip_address character varying(50) DEFAULT ''::character varying,
    user_agent character varying(512) DEFAULT ''::character varying,
    location character varying(512) DEFAULT ''::character varying,
    nid uuid NOT NULL,
    session_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.session_devices OWNER TO test;

--
-- Name: session_token_exchanges; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.session_token_exchanges (
    id uuid NOT NULL,
    nid uuid NOT NULL,
    flow_id uuid NOT NULL,
    session_id uuid,
    init_code character varying(64) NOT NULL,
    return_to_code character varying(64) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.session_token_exchanges OWNER TO test;

--
-- Name: sessions; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.sessions (
    id uuid NOT NULL,
    issued_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    authenticated_at timestamp without time zone NOT NULL,
    identity_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    token character varying(39),
    active boolean DEFAULT false,
    nid uuid,
    logout_token character varying(39),
    aal character varying(4) DEFAULT 'aal1'::character varying NOT NULL,
    authentication_methods jsonb NOT NULL
);


ALTER TABLE public.sessions OWNER TO test;

--
-- Data for Name: continuity_containers; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.continuity_containers (id, identity_id, name, payload, expires_at, created_at, updated_at, nid) FROM stdin;
\.


--
-- Data for Name: courier_message_dispatches; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.courier_message_dispatches (id, message_id, status, error, nid, created_at, updated_at) FROM stdin;
601011e2-f16b-4f2c-91b8-a47915b1fcfb	a720fbb9-2da5-4231-a73a-253bb275f4a0	success	null	8712f917-6eee-4fc0-98b3-364070db5d96	2024-02-09 16:00:36.612791	2024-02-09 16:00:36.612791
\.


--
-- Data for Name: courier_messages; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.courier_messages (id, type, status, body, subject, recipient, created_at, updated_at, template_type, template_data, nid, send_count) FROM stdin;
a720fbb9-2da5-4231-a73a-253bb275f4a0	1	2	Hi,\n\nplease verify your account by entering the following code:\n\n143443\n\nor clicking the following link:\n\nhttps://juicer-dev.xyz/kratos/self-service/verification?code=143443&flow=b7f3f3c0-334a-4437-bbb9-2e81152dca1b\n	Please verify your email address	test@test.com	2024-02-09 16:00:35.783748	2024-02-09 16:00:35.783748	verification_code_valid	\\x7b22546f223a227465737440746573742e636f6d222c22566572696669636174696f6e55524c223a2268747470733a2f2f6a75696365722d6465762e78797a2f6b7261746f732f73656c662d736572766963652f766572696669636174696f6e3f636f64653d3134333434335c7530303236666c6f773d62376633663363302d333334612d343433372d626262392d326538313135326463613162222c22566572696669636174696f6e436f6465223a22313433343433222c224964656e74697479223a7b22637265617465645f6174223a22323032342d30322d30395431363a30303a33352e3736333838375a222c226964223a2232393835623062612d643739362d343237632d623161652d626166623236636232356630222c226d657461646174615f7075626c6963223a6e756c6c2c227265636f766572795f616464726573736573223a5b7b22637265617465645f6174223a22323032342d30322d30395431363a30303a33352e3736363939345a222c226964223a2265303339343937342d366432342d343933312d383130352d633364633533383366613932222c22757064617465645f6174223a22323032342d30322d30395431363a30303a33352e3736363939345a222c2276616c7565223a227465737440746573742e636f6d222c22766961223a22656d61696c227d5d2c22736368656d615f6964223a22637573746f6d6572222c22736368656d615f75726c223a2268747470733a2f2f6a75696365722d6465762e78797a2f6b7261746f732f736368656d61732f5933567a644739745a5849222c227374617465223a22616374697665222c2273746174655f6368616e6765645f6174223a22323032342d30322d30395431363a30303a33352e3736323931323230345a222c22747261697473223a7b226176617461725f75726c223a22222c22656d61696c223a227465737440746573742e636f6d222c2266697273745f6e616d65223a22546573744669727374222c226c6173745f6e616d65223a22546573744c617374227d2c22757064617465645f6174223a22323032342d30322d30395431363a30303a33352e3736333838375a222c2276657269666961626c655f616464726573736573223a5b7b22637265617465645f6174223a22323032342d30322d30395431363a30303a33352e3736353136395a222c226964223a2266356130323266622d376664612d346434322d623330352d363936663335356437373263222c22737461747573223a2270656e64696e67222c22757064617465645f6174223a22323032342d30322d30395431363a30303a33352e3736353136395a222c2276616c7565223a227465737440746573742e636f6d222c227665726966696564223a66616c73652c22766961223a22656d61696c227d5d7d7d	8712f917-6eee-4fc0-98b3-364070db5d96	1
\.


--
-- Data for Name: identities; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.identities (id, schema_id, traits, created_at, updated_at, nid, state, state_changed_at, metadata_public, metadata_admin) FROM stdin;
2985b0ba-d796-427c-b1ae-bafb26cb25f0	customer	{"email": "test@test.com", "last_name": "TestLast", "avatar_url": "", "first_name": "TestFirst"}	2024-02-09 16:00:35.763887	2024-02-09 16:00:35.763887	8712f917-6eee-4fc0-98b3-364070db5d96	active	2024-02-09 16:00:35.762912	\N	\N
\.


--
-- Data for Name: identity_credential_identifiers; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.identity_credential_identifiers (id, identifier, identity_credential_id, created_at, updated_at, nid, identity_credential_type_id) FROM stdin;
f6539ad8-0b7b-4618-8e58-6d29fc69dc48	test@test.com	2b137dd9-94d4-48af-bcdb-5db52965c94e	2024-02-09 16:00:35.771996	2024-02-09 16:00:35.771996	8712f917-6eee-4fc0-98b3-364070db5d96	78c1b41d-8341-4507-aa60-aff1d4369670
\.


--
-- Data for Name: identity_credential_types; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.identity_credential_types (id, name) FROM stdin;
78c1b41d-8341-4507-aa60-aff1d4369670	password
6fa5e2e0-bfce-4631-b62b-cf2b0252b289	oidc
5e29b036-aa47-457f-9fe6-aa8b854a752b	totp
567a0730-7f48-4dd7-a13d-df87a51c245f	lookup_secret
6b213fa0-e6ad-46cb-8878-b088d2ce2e3c	webauthn
\.


--
-- Data for Name: identity_credentials; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.identity_credentials (id, config, identity_credential_type_id, identity_id, created_at, updated_at, nid, version) FROM stdin;
2b137dd9-94d4-48af-bcdb-5db52965c94e	{"hashed_password": "$2a$12$vbCcieZZR7kaZyrPJZ21X./dPZ5dPSw15wxkX2pXhSBRinf1AmEla"}	78c1b41d-8341-4507-aa60-aff1d4369670	2985b0ba-d796-427c-b1ae-bafb26cb25f0	2024-02-09 16:00:35.770638	2024-02-09 16:00:35.770638	8712f917-6eee-4fc0-98b3-364070db5d96	0
\.


--
-- Data for Name: identity_recovery_addresses; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.identity_recovery_addresses (id, via, value, identity_id, created_at, updated_at, nid) FROM stdin;
e0394974-6d24-4931-8105-c3dc5383fa92	email	test@test.com	2985b0ba-d796-427c-b1ae-bafb26cb25f0	2024-02-09 16:00:35.766994	2024-02-09 16:00:35.766994	8712f917-6eee-4fc0-98b3-364070db5d96
\.


--
-- Data for Name: identity_recovery_codes; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.identity_recovery_codes (id, code, used_at, identity_recovery_address_id, code_type, expires_at, issued_at, selfservice_recovery_flow_id, created_at, updated_at, nid, identity_id) FROM stdin;
\.


--
-- Data for Name: identity_recovery_tokens; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.identity_recovery_tokens (id, token, used, used_at, identity_recovery_address_id, selfservice_recovery_flow_id, created_at, updated_at, expires_at, issued_at, nid, identity_id, token_type) FROM stdin;
\.


--
-- Data for Name: identity_verifiable_addresses; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.identity_verifiable_addresses (id, status, via, verified, value, verified_at, identity_id, created_at, updated_at, nid) FROM stdin;
f5a022fb-7fda-4d42-b305-696f355d772c	completed	email	t	test@test.com	2024-02-09 16:00:43.544788	2985b0ba-d796-427c-b1ae-bafb26cb25f0	2024-02-09 16:00:35.765169	2024-02-09 16:00:35.765169	8712f917-6eee-4fc0-98b3-364070db5d96
\.


--
-- Data for Name: identity_verification_codes; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.identity_verification_codes (id, code_hmac, used_at, identity_verifiable_address_id, expires_at, issued_at, selfservice_verification_flow_id, created_at, updated_at, nid) FROM stdin;
f2c05a0b-2735-4c51-a2f8-31db02d3c69e	b421331b2599bf1d15b3dce1f759353f2aa52f4fbbb613ccd6ea7fc25fd94bd2	2024-02-09 16:00:43.53805	f5a022fb-7fda-4d42-b305-696f355d772c	2024-02-09 16:15:35.779781	2024-02-09 16:00:35.779781	b7f3f3c0-334a-4437-bbb9-2e81152dca1b	2024-02-09 16:00:35.779879	2024-02-09 16:00:35.779879	8712f917-6eee-4fc0-98b3-364070db5d96
\.


--
-- Data for Name: identity_verification_tokens; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.identity_verification_tokens (id, token, used, used_at, expires_at, issued_at, identity_verifiable_address_id, selfservice_verification_flow_id, created_at, updated_at, nid) FROM stdin;
\.


--
-- Data for Name: keto_relation_tuples; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.keto_relation_tuples (shard_id, nid, namespace, object, relation, subject_id, subject_set_namespace, subject_set_object, subject_set_relation, commit_time) FROM stdin;
\.


--
-- Data for Name: keto_uuid_mappings; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.keto_uuid_mappings (id, string_representation) FROM stdin;
\.


--
-- Data for Name: networks; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.networks (id, created_at, updated_at) FROM stdin;
8712f917-6eee-4fc0-98b3-364070db5d96	2024-02-09 15:58:32.663115	2024-02-09 15:58:32.663115
\.


--
-- Data for Name: schema_migration; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.schema_migration (version, version_self) FROM stdin;
20150100000001000000	0
20201110175414000000	0
20201110175414000001	0
20210623162417000000	0
20210623162417000001	0
20210623162417000002	0
20210623162417000003	0
20210914134624000000	0
20220217152313000000	0
20220512151000000000	0
20220513200300000000	0
20220513200400000000	0
20220513200400000001	0
20220513200500000000	0
20220513200600000000	0
20220513200600000001	0
20230228091200000000	0
20191100000001000000	0
20191100000001000001	0
20191100000001000002	0
20191100000001000003	0
20191100000001000004	0
20191100000001000005	0
20191100000002000000	0
20191100000002000001	0
20191100000002000002	0
20191100000002000003	0
20191100000002000004	0
20191100000003000000	0
20191100000004000000	0
20191100000006000000	0
20191100000007000000	0
20191100000008000000	0
20191100000008000001	0
20191100000008000002	0
20191100000008000003	0
20191100000008000004	0
20191100000008000005	0
20191100000010000000	0
20191100000010000001	0
20191100000011000000	0
20191100000012000000	0
20200317160354000000	0
20200317160354000001	0
20200317160354000002	0
20200317160354000003	0
20200317160354000004	0
20200401183443000000	0
20200402142539000000	0
20200402142539000001	0
20200402142539000002	0
20200519101057000000	0
20200519101057000001	0
20200519101057000002	0
20200519101057000003	0
20200519101057000004	0
20200519101057000005	0
20200519101057000006	0
20200519101057000007	0
20200601101000000000	0
20200605111551000000	0
20200605111551000001	0
20200605111551000002	0
20200607165100000000	0
20200607165100000001	0
20200705105359000000	0
20200810141652000000	0
20200810141652000001	0
20200810141652000002	0
20200810141652000003	0
20200810141652000004	0
20200810161022000000	0
20200810161022000001	0
20200810161022000002	0
20200810161022000003	0
20200810161022000004	0
20200810161022000005	0
20200810161022000006	0
20200810161022000007	0
20200810161022000008	0
20200810162450000000	0
20200810162450000001	0
20200810162450000002	0
20200810162450000003	0
20200812124254000000	0
20200812124254000001	0
20200812124254000002	0
20200812124254000003	0
20200812124254000004	0
20200812160551000000	0
20200830121710000000	0
20200830130642000000	0
20200830130642000001	0
20200830130642000002	0
20200830130642000003	0
20200830130642000004	0
20200830130642000005	0
20200830130642000006	0
20200830130642000007	0
20200830130643000000	0
20200830130644000000	0
20200830130644000001	0
20200830130645000000	0
20200830130646000000	0
20200830130646000001	0
20200830130646000002	0
20200830154602000000	0
20200830154602000001	0
20200830154602000002	0
20200830154602000003	0
20200830154602000004	0
20200830172221000000	0
20200830172221000001	0
20200830172221000002	0
20200830172221000003	0
20200831110752000000	0
20200831110752000001	0
20200831110752000002	0
20200831110752000003	0
20200831110752000004	0
20200831110752000005	0
20200831110752000006	0
20200831110752000007	0
20201201161451000000	0
20201201161451000001	0
20210307130558000000	0
20210307130559000000	0
20210307130559000001	0
20210311102338000000	0
20210311102338000001	0
20210311102338000002	0
20210311102338000003	0
20210311102338000004	0
20210311102338000005	0
20210311102338000006	0
20210311102338000007	0
20210311102338000008	0
20210311102338000009	0
20210311102338000010	0
20210311102338000011	0
20210311102338000012	0
20210311102338000013	0
20210311102338000014	0
20210311102338000015	0
20210311102338000016	0
20210311102338000017	0
20210311102338000018	0
20210311102338000019	0
20210311102338000020	0
20210311102338000021	0
20210311102338000022	0
20210311102338000023	0
20210311102338000024	0
20210410175418000000	0
20210410175418000001	0
20210410175418000002	0
20210410175418000003	0
20210410175418000004	0
20210410175418000005	0
20210410175418000006	0
20210410175418000007	0
20210410175418000008	0
20210410175418000009	0
20210410175418000010	0
20210410175418000011	0
20210410175418000012	0
20210410175418000013	0
20210410175418000014	0
20210410175418000015	0
20210410175418000016	0
20210410175418000017	0
20210410175418000018	0
20210410175418000019	0
20210410175418000020	0
20210410175418000021	0
20210410175418000022	0
20210410175418000023	0
20210410175418000024	0
20210410175418000025	0
20210410175418000026	0
20210410175418000027	0
20210410175418000028	0
20210410175418000029	0
20210410175418000030	0
20210410175418000031	0
20210410175418000032	0
20210410175418000033	0
20210410175418000034	0
20210410175418000035	0
20210410175418000036	0
20210410175418000037	0
20210410175418000038	0
20210410175418000039	0
20210410175418000040	0
20210410175418000041	0
20210410175418000042	0
20210410175418000043	0
20210410175418000044	0
20210410175418000045	0
20210410175418000046	0
20210410175418000047	0
20210410175418000048	0
20210410175418000049	0
20210410175418000050	0
20210410175418000051	0
20210410175418000052	0
20210410175418000053	0
20210410175418000054	0
20210410175418000055	0
20210410175418000056	0
20210410175418000057	0
20210410175418000058	0
20210410175418000059	0
20210410175418000060	0
20210410175418000061	0
20210410175418000062	0
20210410175418000063	0
20210410175418000064	0
20210410175418000065	0
20210410175418000066	0
20210410175418000067	0
20210410175418000068	0
20210410175418000069	0
20210410175418000070	0
20210410175418000071	0
20210410175418000072	0
20210410175418000073	0
20210410175418000074	0
20210410175418000075	0
20210410175418000076	0
20210410175418000077	0
20210410175418000078	0
20210410175418000079	0
20210410175418000080	0
20210410175418000081	0
20210410175418000082	0
20210410175418000083	0
20210410175418000084	0
20210410175418000085	0
20210410175418000086	0
20210410175418000087	0
20210410175418000088	0
20210410175418000089	0
20210504121624000000	0
20210504121624000001	0
20210618103120000000	0
20210618103120000001	0
20210618103120000002	0
20210618103120000003	0
20210618103120000004	0
20210805112414000000	0
20210805112414000001	0
20210805112414000002	0
20210805122535000000	0
20210810153530000000	0
20210810153530000001	0
20210810153530000002	0
20210810153530000003	0
20210810153530000004	0
20210813150152000000	0
20210816113956000000	0
20210816142650000000	0
20210816142650000001	0
20210816142650000002	0
20210816142650000003	0
20210816142650000004	0
20210816142650000005	0
20210817181232000000	0
20210817181232000001	0
20210817181232000002	0
20210817181232000003	0
20210817181232000004	0
20210817181232000005	0
20210829131458000000	0
20210913095309000000	0
20210913095309000001	0
20210913095309000002	0
20210913095309000003	0
20210913095309000004	0
20220118104539000000	0
20220118104539000001	0
20220118104539000002	0
20220118104539000003	0
20220301102701000000	0
20220301102701000001	0
20220420102701000000	0
20220512102703000000	0
20220607000001000000	0
20220610155809000000	0
20220802103909000000	0
20220824165300000000	0
20220824165300000001	0
20220824165300000002	0
20220825134336000000	0
20220825134336000001	0
20220901123209000000	0
20220907132836000000	0
20221024182336000000	0
20221205092803000000	0
20221214101328000000	0
20221220124639000000	0
20230104193739000000	0
20230216142104000000	0
20230313141439000000	0
20230313141439000001	0
20230322144139000001	0
20230405000000000001	0
20230614000001000000	0
20230619000000000001	0
20230626000000000001	0
20230705000000000001	0
\.


--
-- Data for Name: selfservice_errors; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.selfservice_errors (id, errors, seen_at, was_seen, created_at, updated_at, csrf_token, nid) FROM stdin;
\.


--
-- Data for Name: selfservice_login_flows; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.selfservice_login_flows (id, request_url, issued_at, expires_at, active_method, csrf_token, created_at, updated_at, forced, type, ui, nid, requested_aal, internal_context, oauth2_login_challenge, oauth2_login_challenge_data) FROM stdin;
9e2e6257-289e-4607-954d-afbc4354cafb	https://juicer-dev.xyz/self-service/login/browser	2024-02-09 16:00:03.902558	2024-02-09 16:10:03.902558		ySvK7RCGhAVz6tAkxWmTQYCpD9wdYD3dfe9T881spRsPRRiVUu/H1QEMCoPE4PfvcZTlkNs60qLIYGmYk9iCkg==	2024-02-09 16:00:03.907012	2024-02-09 16:00:03.907012	f	browser	{"nodes": [{"meta": {"label": {"id": 1010002, "text": "Sign in with discord", "type": "info", "context": {"provider": "discord"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "discord", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with facebook", "type": "info", "context": {"provider": "facebook"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "facebook", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with github", "type": "info", "context": {"provider": "github"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "github", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with google", "type": "info", "context": {"provider": "google"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "google", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with slack", "type": "info", "context": {"provider": "slack"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "slack", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with spotify", "type": "info", "context": {"provider": "spotify"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "spotify", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with twitch", "type": "info", "context": {"provider": "twitch"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "twitch", "disabled": false, "node_type": "input"}}, {"meta": {}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "csrf_token", "type": "hidden", "value": "ySvK7RCGhAVz6tAkxWmTQYCpD9wdYD3dfe9T881spRsPRRiVUu/H1QEMCoPE4PfvcZTlkNs60qLIYGmYk9iCkg==", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070004, "text": "ID", "type": "info"}}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "identifier", "type": "text", "value": "", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070001, "text": "Password", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "password", "type": "password", "disabled": false, "required": true, "node_type": "input", "autocomplete": "current-password"}}, {"meta": {"label": {"id": 1010001, "text": "Sign in", "type": "info", "context": {}}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "method", "type": "submit", "value": "password", "disabled": false, "node_type": "input"}}], "action": "https://juicer-dev.xyz/kratos/self-service/login?flow=9e2e6257-289e-4607-954d-afbc4354cafb", "method": "POST"}	8712f917-6eee-4fc0-98b3-364070db5d96	aal1	{}	\N	\N
f28668ba-e0ab-4bc2-8356-363f413116a2	https://juicer-dev.xyz/self-service/login/browser	2024-02-09 16:00:44.781616	2024-02-09 16:10:44.781616		f79fViPZEIRWFclYVJX9tRS+/6dOIVsS/cLsCXrbsam9hYuF2zwGJDgjWrZAQu0zPDplTkIJR8dSzN2HU/7N0Q==	2024-02-09 16:00:44.785913	2024-02-09 16:00:44.785913	f	browser	{"nodes": [{"meta": {"label": {"id": 1010002, "text": "Sign in with discord", "type": "info", "context": {"provider": "discord"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "discord", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with facebook", "type": "info", "context": {"provider": "facebook"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "facebook", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with github", "type": "info", "context": {"provider": "github"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "github", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with google", "type": "info", "context": {"provider": "google"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "google", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with slack", "type": "info", "context": {"provider": "slack"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "slack", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with spotify", "type": "info", "context": {"provider": "spotify"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "spotify", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with twitch", "type": "info", "context": {"provider": "twitch"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "twitch", "disabled": false, "node_type": "input"}}, {"meta": {}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "csrf_token", "type": "hidden", "value": "f79fViPZEIRWFclYVJX9tRS+/6dOIVsS/cLsCXrbsam9hYuF2zwGJDgjWrZAQu0zPDplTkIJR8dSzN2HU/7N0Q==", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070004, "text": "ID", "type": "info"}}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "identifier", "type": "text", "value": "", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070001, "text": "Password", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "password", "type": "password", "disabled": false, "required": true, "node_type": "input", "autocomplete": "current-password"}}, {"meta": {"label": {"id": 1010001, "text": "Sign in", "type": "info", "context": {}}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "method", "type": "submit", "value": "password", "disabled": false, "node_type": "input"}}], "action": "https://juicer-dev.xyz/kratos/self-service/login?flow=f28668ba-e0ab-4bc2-8356-363f413116a2", "method": "POST"}	8712f917-6eee-4fc0-98b3-364070db5d96	aal1	{}	\N	\N
262bbac1-e147-4d80-8f4e-ed57e083b1af	https://juicer-dev.xyz/self-service/login/browser	2024-02-09 16:00:48.551613	2024-02-09 16:10:48.551613	password	pF1uX2svH76q18gL8Zw+tGvgCKTOvPGHDh/HTyrRrlhmZ7qMk8oJHsThW+XlSy4yQ2SSTcKU7VKhEfbBA/TSIA==	2024-02-09 16:00:48.557903	2024-02-09 16:00:48.557903	f	browser	{"nodes": [{"meta": {"label": {"id": 1010002, "text": "Sign in with discord", "type": "info", "context": {"provider": "discord"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "discord", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with facebook", "type": "info", "context": {"provider": "facebook"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "facebook", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with github", "type": "info", "context": {"provider": "github"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "github", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with google", "type": "info", "context": {"provider": "google"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "google", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with slack", "type": "info", "context": {"provider": "slack"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "slack", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with spotify", "type": "info", "context": {"provider": "spotify"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "spotify", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with twitch", "type": "info", "context": {"provider": "twitch"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "twitch", "disabled": false, "node_type": "input"}}, {"meta": {}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "csrf_token", "type": "hidden", "value": "pF1uX2svH76q18gL8Zw+tGvgCKTOvPGHDh/HTyrRrlhmZ7qMk8oJHsThW+XlSy4yQ2SSTcKU7VKhEfbBA/TSIA==", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070004, "text": "ID", "type": "info"}}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "identifier", "type": "text", "value": "", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070001, "text": "Password", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "password", "type": "password", "disabled": false, "required": true, "node_type": "input", "autocomplete": "current-password"}}, {"meta": {"label": {"id": 1010001, "text": "Sign in", "type": "info", "context": {}}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "method", "type": "submit", "value": "password", "disabled": false, "node_type": "input"}}], "action": "https://juicer-dev.xyz/kratos/self-service/login?flow=262bbac1-e147-4d80-8f4e-ed57e083b1af", "method": "POST"}	8712f917-6eee-4fc0-98b3-364070db5d96	aal1	{}	\N	\N
\.


--
-- Data for Name: selfservice_recovery_flows; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.selfservice_recovery_flows (id, request_url, issued_at, expires_at, active_method, csrf_token, state, recovered_identity_id, created_at, updated_at, type, ui, nid, submit_count, skip_csrf_check) FROM stdin;
d9172e9f-9356-4b21-a5b1-06da202bffd7	https://juicer-dev.xyz/self-service/recovery/browser	2024-02-09 16:00:07.282344	2024-02-09 17:00:07.282344	code	b7XlqUsIWBTQFNGk/A6cnngeferVwv0vyUGQaTWTq0Gp2zfRCWEbxKLyCwP9h/gwiSOXphOYElB8zqoCayeMyA==	choose_method	\N	2024-02-09 16:00:07.282464	2024-02-09 16:00:07.282464	browser	{"nodes": [{"meta": {}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "csrf_token", "type": "hidden", "value": "b7XlqUsIWBTQFNGk/A6cnngeferVwv0vyUGQaTWTq0Gp2zfRCWEbxKLyCwP9h/gwiSOXphOYElB8zqoCayeMyA==", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070007, "text": "Email", "type": "info"}}, "type": "input", "group": "code", "messages": [], "attributes": {"name": "email", "type": "email", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070005, "text": "Submit", "type": "info"}}, "type": "input", "group": "code", "messages": [], "attributes": {"name": "method", "type": "submit", "value": "code", "disabled": false, "node_type": "input"}}], "action": "https://juicer-dev.xyz/kratos/self-service/recovery?flow=d9172e9f-9356-4b21-a5b1-06da202bffd7", "method": "POST"}	8712f917-6eee-4fc0-98b3-364070db5d96	0	f
0b1a355d-f1ff-4054-b3da-a51bf46678ee	https://juicer-dev.xyz/self-service/recovery/browser	2024-02-09 16:00:07.330801	2024-02-09 17:00:07.330801	code	BTrVDo2rHgjWlFkfugIEEKpmjHbdAeholDVsqWDqj4nDVAd2z8Jd2KRyg7i7i2C+W1tmOhtbBxchulbCPl6oAA==	choose_method	\N	2024-02-09 16:00:07.330926	2024-02-09 16:00:07.330926	browser	{"nodes": [{"meta": {}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "csrf_token", "type": "hidden", "value": "BTrVDo2rHgjWlFkfugIEEKpmjHbdAeholDVsqWDqj4nDVAd2z8Jd2KRyg7i7i2C+W1tmOhtbBxchulbCPl6oAA==", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070007, "text": "Email", "type": "info"}}, "type": "input", "group": "code", "messages": [], "attributes": {"name": "email", "type": "email", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070005, "text": "Submit", "type": "info"}}, "type": "input", "group": "code", "messages": [], "attributes": {"name": "method", "type": "submit", "value": "code", "disabled": false, "node_type": "input"}}], "action": "https://juicer-dev.xyz/kratos/self-service/recovery?flow=0b1a355d-f1ff-4054-b3da-a51bf46678ee", "method": "POST"}	8712f917-6eee-4fc0-98b3-364070db5d96	0	f
9ddb6a2a-e393-40eb-952d-dc8b5f65a47c	https://juicer-dev.xyz/self-service/recovery/browser	2024-02-09 16:00:07.332264	2024-02-09 17:00:07.332264	code	hCafcubfWv9NJwm3P0is6i9rH4AXA8JgJu6vJBSqTUtCSE0KpLYZLz/B0xA+wchE3lb1zNFZLR+TYZVPSh5qwg==	choose_method	\N	2024-02-09 16:00:07.33235	2024-02-09 16:00:07.33235	browser	{"nodes": [{"meta": {}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "csrf_token", "type": "hidden", "value": "hCafcubfWv9NJwm3P0is6i9rH4AXA8JgJu6vJBSqTUtCSE0KpLYZLz/B0xA+wchE3lb1zNFZLR+TYZVPSh5qwg==", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070007, "text": "Email", "type": "info"}}, "type": "input", "group": "code", "messages": [], "attributes": {"name": "email", "type": "email", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070005, "text": "Submit", "type": "info"}}, "type": "input", "group": "code", "messages": [], "attributes": {"name": "method", "type": "submit", "value": "code", "disabled": false, "node_type": "input"}}], "action": "https://juicer-dev.xyz/kratos/self-service/recovery?flow=9ddb6a2a-e393-40eb-952d-dc8b5f65a47c", "method": "POST"}	8712f917-6eee-4fc0-98b3-364070db5d96	0	f
22f1a3d5-3e55-428f-aa16-ca55637e6050	https://juicer-dev.xyz/self-service/recovery/browser	2024-02-09 16:00:13.304148	2024-02-09 17:00:13.304148	code	g+5I0FL2pEUeYMMxYkPkp4rIALCv9wk9+GcVVBMomYRFgJqoEJ/nlWyGGZZjyoAJe/Xq/Gmt5kJN6C8/TZy+DQ==	choose_method	\N	2024-02-09 16:00:13.304293	2024-02-09 16:00:13.304293	browser	{"nodes": [{"meta": {}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "csrf_token", "type": "hidden", "value": "g+5I0FL2pEUeYMMxYkPkp4rIALCv9wk9+GcVVBMomYRFgJqoEJ/nlWyGGZZjyoAJe/Xq/Gmt5kJN6C8/TZy+DQ==", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070007, "text": "Email", "type": "info"}}, "type": "input", "group": "code", "messages": [], "attributes": {"name": "email", "type": "email", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070005, "text": "Submit", "type": "info"}}, "type": "input", "group": "code", "messages": [], "attributes": {"name": "method", "type": "submit", "value": "code", "disabled": false, "node_type": "input"}}], "action": "https://juicer-dev.xyz/kratos/self-service/recovery?flow=22f1a3d5-3e55-428f-aa16-ca55637e6050", "method": "POST"}	8712f917-6eee-4fc0-98b3-364070db5d96	0	f
\.


--
-- Data for Name: selfservice_registration_flows; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.selfservice_registration_flows (id, request_url, issued_at, expires_at, active_method, csrf_token, created_at, updated_at, type, ui, nid, internal_context, oauth2_login_challenge, oauth2_login_challenge_data) FROM stdin;
7a1bc20a-9c00-4d86-aaf2-c3b884fe9ea9	https://juicer-dev.xyz/self-service/registration/browser	2024-02-09 16:00:02.204137	2024-02-09 16:10:02.204137		O3cYMh5E89+zN+BTqMpl3k5mVq0wK55pJyLXvHyLPRX9GcpKXC2wD8HROvSpQwFwv1u84fZxcRaSre3XIj8anA==	2024-02-09 16:00:02.208231	2024-02-09 16:00:02.208231	browser	{"nodes": [{"meta": {}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "csrf_token", "type": "hidden", "value": "O3cYMh5E89+zN+BTqMpl3k5mVq0wK55pJyLXvHyLPRX9GcpKXC2wD8HROvSpQwFwv1u84fZxcRaSre3XIj8anA==", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with discord", "type": "info", "context": {"provider": "discord"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "discord", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with facebook", "type": "info", "context": {"provider": "facebook"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "facebook", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with github", "type": "info", "context": {"provider": "github"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "github", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with google", "type": "info", "context": {"provider": "google"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "google", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with slack", "type": "info", "context": {"provider": "slack"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "slack", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with spotify", "type": "info", "context": {"provider": "spotify"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "spotify", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with twitch", "type": "info", "context": {"provider": "twitch"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "twitch", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "E-Mail", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "traits.email", "type": "email", "disabled": false, "required": true, "node_type": "input", "autocomplete": "email"}}, {"meta": {"label": {"id": 1070001, "text": "Password", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "password", "type": "password", "disabled": false, "required": true, "node_type": "input", "autocomplete": "new-password"}}, {"meta": {"label": {"id": 1070002, "text": "First name", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "traits.first_name", "type": "text", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "Last name", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "traits.last_name", "type": "text", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "Picture", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "traits.avatar_url", "type": "text", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040001, "text": "Sign up", "type": "info", "context": {}}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "method", "type": "submit", "value": "password", "disabled": false, "node_type": "input"}}], "action": "https://juicer-dev.xyz/kratos/self-service/registration?flow=7a1bc20a-9c00-4d86-aaf2-c3b884fe9ea9", "method": "POST"}	8712f917-6eee-4fc0-98b3-364070db5d96	{}	\N	\N
0f5a8f69-779e-4ab8-ad35-dbfcea24d9ed	https://juicer-dev.xyz/self-service/registration/browser	2024-02-09 16:00:17.350746	2024-02-09 16:10:17.350746		hkzv0INmFDG/4KB3CwmTjhElRvFY72Hmkh/T5Zs3HmxAIj2owQ9X4c0GetAKgPcg4BisvZ61jpknkOmOxYM55Q==	2024-02-09 16:00:17.355019	2024-02-09 16:00:17.355019	browser	{"nodes": [{"meta": {}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "csrf_token", "type": "hidden", "value": "hkzv0INmFDG/4KB3CwmTjhElRvFY72Hmkh/T5Zs3HmxAIj2owQ9X4c0GetAKgPcg4BisvZ61jpknkOmOxYM55Q==", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with discord", "type": "info", "context": {"provider": "discord"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "discord", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with facebook", "type": "info", "context": {"provider": "facebook"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "facebook", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with github", "type": "info", "context": {"provider": "github"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "github", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with google", "type": "info", "context": {"provider": "google"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "google", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with slack", "type": "info", "context": {"provider": "slack"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "slack", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with spotify", "type": "info", "context": {"provider": "spotify"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "spotify", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with twitch", "type": "info", "context": {"provider": "twitch"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "twitch", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "E-Mail", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "traits.email", "type": "email", "disabled": false, "required": true, "node_type": "input", "autocomplete": "email"}}, {"meta": {"label": {"id": 1070001, "text": "Password", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "password", "type": "password", "disabled": false, "required": true, "node_type": "input", "autocomplete": "new-password"}}, {"meta": {"label": {"id": 1070002, "text": "First name", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "traits.first_name", "type": "text", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "Last name", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "traits.last_name", "type": "text", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "Picture", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "traits.avatar_url", "type": "text", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040001, "text": "Sign up", "type": "info", "context": {}}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "method", "type": "submit", "value": "password", "disabled": false, "node_type": "input"}}], "action": "https://juicer-dev.xyz/kratos/self-service/registration?flow=0f5a8f69-779e-4ab8-ad35-dbfcea24d9ed", "method": "POST"}	8712f917-6eee-4fc0-98b3-364070db5d96	{}	\N	\N
87df423a-256b-4c03-8aef-03bcbb0b1dcb	https://juicer-dev.xyz/self-service/registration/browser	2024-02-09 16:00:35.459146	2024-02-09 16:10:35.459146		Zml8MM1Ze7tQRmfFrno6Ces29CKARgvZ4GRIA/7kczOgB65IjzA4ayKgvWKv816nGgsebkYc5KZV63JooFBUug==	2024-02-09 16:00:35.466022	2024-02-09 16:00:35.466022	browser	{"nodes": [{"meta": {}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "csrf_token", "type": "hidden", "value": "Zml8MM1Ze7tQRmfFrno6Ces29CKARgvZ4GRIA/7kczOgB65IjzA4ayKgvWKv816nGgsebkYc5KZV63JooFBUug==", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with discord", "type": "info", "context": {"provider": "discord"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "discord", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with facebook", "type": "info", "context": {"provider": "facebook"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "facebook", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with github", "type": "info", "context": {"provider": "github"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "github", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with google", "type": "info", "context": {"provider": "google"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "google", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with slack", "type": "info", "context": {"provider": "slack"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "slack", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with spotify", "type": "info", "context": {"provider": "spotify"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "spotify", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with twitch", "type": "info", "context": {"provider": "twitch"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "twitch", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "E-Mail", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "traits.email", "type": "email", "disabled": false, "required": true, "node_type": "input", "autocomplete": "email"}}, {"meta": {"label": {"id": 1070001, "text": "Password", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "password", "type": "password", "disabled": false, "required": true, "node_type": "input", "autocomplete": "new-password"}}, {"meta": {"label": {"id": 1070002, "text": "First name", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "traits.first_name", "type": "text", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "Last name", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "traits.last_name", "type": "text", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "Picture", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "traits.avatar_url", "type": "text", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040001, "text": "Sign up", "type": "info", "context": {}}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "method", "type": "submit", "value": "password", "disabled": false, "node_type": "input"}}], "action": "https://juicer-dev.xyz/kratos/self-service/registration?flow=87df423a-256b-4c03-8aef-03bcbb0b1dcb", "method": "POST"}	8712f917-6eee-4fc0-98b3-364070db5d96	{}	\N	\N
\.


--
-- Data for Name: selfservice_settings_flows; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.selfservice_settings_flows (id, request_url, issued_at, expires_at, identity_id, created_at, updated_at, active_method, state, type, ui, nid, internal_context) FROM stdin;
fffa7987-5d71-4aaa-8ce6-8b39934f0bc0	https://juicer-dev.xyz/self-service/settings/browser	2024-02-09 16:01:02.168319	2024-02-09 17:01:02.168319	2985b0ba-d796-427c-b1ae-bafb26cb25f0	2024-02-09 16:01:02.1986	2024-02-09 16:01:02.1986	\N	show_form	browser	{"nodes": [{"meta": {}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "csrf_token", "type": "hidden", "value": "2TC5uhGdbLOD6nXHsb+rRdcPJ9dyDR2W0Ou137AwQoxDJ2A9yh6dDrudCCs8hw7ZX6VXZwGc4ajXgrzjtnVt5g==", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "E-Mail", "type": "info"}}, "type": "input", "group": "profile", "messages": [], "attributes": {"name": "traits.email", "type": "email", "value": "test@test.com", "disabled": false, "required": true, "node_type": "input", "autocomplete": "email"}}, {"meta": {"label": {"id": 1070002, "text": "First name", "type": "info"}}, "type": "input", "group": "profile", "messages": [], "attributes": {"name": "traits.first_name", "type": "text", "value": "TestFirst", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "Last name", "type": "info"}}, "type": "input", "group": "profile", "messages": [], "attributes": {"name": "traits.last_name", "type": "text", "value": "TestLast", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "Picture", "type": "info"}}, "type": "input", "group": "profile", "messages": [], "attributes": {"name": "traits.avatar_url", "type": "text", "value": "", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1070003, "text": "Save", "type": "info"}}, "type": "input", "group": "profile", "messages": [], "attributes": {"name": "method", "type": "submit", "value": "profile", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1070001, "text": "Password", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "password", "type": "password", "disabled": false, "required": true, "node_type": "input", "autocomplete": "new-password"}}, {"meta": {"label": {"id": 1070003, "text": "Save", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "method", "type": "submit", "value": "password", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1050002, "text": "Link discord", "type": "info", "context": {"provider": "discord"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "link", "type": "submit", "value": "discord", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1050002, "text": "Link facebook", "type": "info", "context": {"provider": "facebook"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "link", "type": "submit", "value": "facebook", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1050002, "text": "Link github", "type": "info", "context": {"provider": "github"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "link", "type": "submit", "value": "github", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1050002, "text": "Link google", "type": "info", "context": {"provider": "google"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "link", "type": "submit", "value": "google", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1050002, "text": "Link slack", "type": "info", "context": {"provider": "slack"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "link", "type": "submit", "value": "slack", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1050002, "text": "Link spotify", "type": "info", "context": {"provider": "spotify"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "link", "type": "submit", "value": "spotify", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1050002, "text": "Link twitch", "type": "info", "context": {"provider": "twitch"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "link", "type": "submit", "value": "twitch", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1050008, "text": "Generate new backup recovery codes", "type": "info"}}, "type": "input", "group": "lookup_secret", "messages": [], "attributes": {"name": "lookup_secret_regenerate", "type": "submit", "value": "true", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1050005, "text": "Authenticator app QR code", "type": "info"}}, "type": "img", "group": "totp", "messages": [], "attributes": {"id": "totp_qr", "src": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAQAAAAEAEAAAAAApiSv5AAAHM0lEQVR4nOyd0W4juQ4Fby7m/3959i1tYAk1aR63saiqp8HYrVaSgkBQIvXn79//CZj/f3sC8l0UAI4CwFEAOH+uf/78JAasgsrzyN0w9BplE7hufsrrvee5VJ92/687yobXt7kCwFEAOAoARwHg/Kn+cx5kdcO86nvpgKr73up7Fdez3bnMZzoP7lJ/I1cAOAoARwHgKACcMgi8yARPm1zfPIA8B23Vs5s5n8er6P7WMu+9+zlcAeAoABwFgKMAcG6CwAzzPOEmoOoGfN1t2XR4WdH9OfK4AsBRADgKAEcB4DwSBF5sNmjPT2QyaPNzeOfsZfe9F08X6rgCwFEAOAoARwHg3ASBmZCkm+fqngk8fzovwajGyzwxH2X+G9/9jVwB4CgAHAWAowBwyiAwU4N6HnketGWqg89zudjMr0umYGaHKwAcBYCjAHAUAM5LEPjERuQ8ZHoiBOuyCby+Nec7XAHgKAAcBYCjAHDKPoGZMGW+fZs5kXdmE3xWn3bzifMZnGfVxepgOaIAcBQAjgLACRWGpDN855HnAeTmRN68tjk9yqa2+Q5XADgKAEcB4CgAnHI7eN6gZdMsuiJTCdwNo84jp0s65i1x5j9HRT1TVwA4CgBHAeAoAJyba+M2gc352XkxRkW3oOJzubl5t8Ezm+Y4ZgJljALAUQA4CgDnJ3O77UUm3NqwKSD5XCOXzLV2mYphM4HyiwLAUQA4CgBn1SJm09MuHfClR9ncCdzNgWYqn3cnBl0B4CgAHAWAowBwyiBwXtd7ZlPEsCmZmAdjFd185yY3N+8TeH7v+R1mAuUXBYCjAHAUAE67T+CmyGIT+lXMZ5AOo+aBYXoben6C00ygFCgAHAWAowBwbs4EtoeJnHjrjrdhPtNNgPa5nGpmLq4AeBQAjgLAUQA4ZRBYkW5Tshk53XPv/Lbze9N9Arvv3TzrdrD8ogBwFACOAsBpVwdfPJHdeiJfd/5ed36bbegz6ebT1aeuAHgUAI4CwFEAODeZwNS9FMcpRILF9Ew31cvpmWbOBFoYIgUKAEcB4CgAnJvq4M9tcXbpvjd9x8gTHQM3jaEv5g1uXnEFgKMAcBQAjgLAafcJ7IcVPTbFIrumKInx5sxnkOmpaGGIHFEAOAoARwHgvHFjSKbN88V8Q7r69CIz502gmemfmP49V+91BcCjAHAUAI4CwLnJBGbOoM3Z3Lmx6Qk4b9qcqZrOzPmdzKwrABwFgKMAcBQAThkEztlcmZYZr3r2TGbO5+AzcydwNV6qu6IrABwFgKMAcBQATigInOegNrm+J5pPZ8Ky81zSN6nM5+cKgEcB4CgAHAWAswoC5xu+mRKHinRxRzekO79jHhhucpvvhJyuAHAUAI4CwFEAODdBYPoquYt0k+rPlYM8W7SRaVddYSZQChQAjgLAUQA4N30Cu3yrCOQ8ynwDedMIpxqlW3G92f7enTt0BYCjAHAUAI4CwPn51l0Vn7sYriKT4ctUAm+ezY/sCgBHAeAoABwFgPNGi5huxutiUzYyryJO33IyD1efaKfTzXx6d7AcUQA4CgBHAeDcZALfGDBypi3TuDrTjHnTyGUz8ub/qrl4JlAKFACOAsBRADgvmcBND7punrAiHQ5utqbnucNdq+bp27oZyP6JRlcAOAoARwHgKACc9t3B87KH87PnT9N1uOmq383296aief7brZ51O1h+UQA4CgBHAeC8BIFP31VxevY8SiZk6r5tfhqyGyJmOh+euQtrXQHgKAAcBYCjAHDKM4GZ5sTdZzP5xCdG7r7t/L2LdNb0HVwB4CgAHAWAowBwfp5oy7IJZ9JtmTeNlzOZu0yInQmJXQHwKAAcBYCjAHBugsCKTDHG+Xvzt6Xric+zOjMPB9O5yH520BUAjgLAUQA4CgCnPBM4DyYyBRVnNgHaPNu4CSoz7bTP43VnWmEmUH5RADgKAEcB4LTvDr44tx9Jt3TJkGk+XX3viVltml5bGCJHFACOAsBRADhln8CKz50O3Lxt07imS7oXYfcdc975C7oCwFEAOAoARwHgxG8MeRn6Sx0DM6Uum+zbmc3vZTOeQaAUKAAcBYCjAHDamcAu/bsq/v29+fm6TUe+Jzoffm4bOvXzugLAUQA4CgBHAeCUZwLn2cF5acX5e+eRN7O6yPyUmRKRajO7IhNYv+IKAEcB4CgAHAWA80ZhyMUmM1Z9mr4U7XxisEu6/U06MKye6OMKAEcB4CgAHAWAcxMEZuhe1bYJdj5XQZtpQ5Mp/Zhz915XADgKAEcB4CgAnEeCwDPdXNo8I9fNr20aUp/n0mUT1s5PSL7iCgBHAeAoABwFgHMTBG5qhzO1vheZnoWbZtGZiuFN2ciZdzauXQHgKAAcBYCjAHDKIDCzYZkJ8zbMyyjmc0nfqVI98clr6FwB4CgAHAWAowBwPtgnUP4LuALAUQA4CgBHAeAoAJx/AgAA//+C02JunQdY/QAAAABJRU5ErkJggg==", "width": 256, "height": 256, "node_type": "img"}}, {"meta": {"label": {"id": 1050017, "text": "This is your authenticator app secret. Use it if you can not scan the QR code.", "type": "info"}}, "type": "text", "group": "totp", "messages": [], "attributes": {"id": "totp_secret_key", "text": {"id": 1050006, "text": "CF4EXMBHQKHH2FH2T6JZHY34JA4BT6OJ", "type": "info", "context": {"secret": "CF4EXMBHQKHH2FH2T6JZHY34JA4BT6OJ"}}, "node_type": "text"}}, {"meta": {"label": {"id": 1070006, "text": "Verify code", "type": "info"}}, "type": "input", "group": "totp", "messages": [], "attributes": {"name": "totp_code", "type": "text", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070003, "text": "Save", "type": "info"}}, "type": "input", "group": "totp", "messages": [], "attributes": {"name": "method", "type": "submit", "value": "totp", "disabled": false, "node_type": "input"}}], "action": "https://juicer-dev.xyz/kratos/self-service/settings?flow=fffa7987-5d71-4aaa-8ce6-8b39934f0bc0", "method": "POST"}	8712f917-6eee-4fc0-98b3-364070db5d96	{"totp_url": "otpauth://totp/juicer-dev.xyz:test@test.com?algorithm=SHA1&digits=6&issuer=juicer-dev.xyz&period=30&secret=CF4EXMBHQKHH2FH2T6JZHY34JA4BT6OJ"}
\.


--
-- Data for Name: selfservice_verification_flows; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.selfservice_verification_flows (id, request_url, issued_at, expires_at, csrf_token, created_at, updated_at, type, state, active_method, ui, nid, submit_count) FROM stdin;
b7f3f3c0-334a-4437-bbb9-2e81152dca1b	https://juicer-dev.xyz/self-service/registration/browser?return_to=	2024-02-09 16:00:35.776842	2024-02-09 17:00:35.776842	F+bT/Cc3ri1QEDAddEyK+TzMy3zauXjjyxS+RPs9yHbV3Acv39K4jT4mo/Ngm5p/FEhRldaRZDZkGo/K0hi0Dg==	2024-02-09 16:00:35.776936	2024-02-09 16:00:35.776936	browser	passed_challenge	code	{"nodes": [{"meta": {}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "csrf_token", "type": "hidden", "value": "F+bT/Cc3ri1QEDAddEyK+TzMy3zauXjjyxS+RPs9yHbV3Acv39K4jT4mo/Ngm5p/FEhRldaRZDZkGo/K0hi0Dg==", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070009, "text": "Continue", "type": "info"}}, "type": "a", "group": "code", "messages": [], "attributes": {"id": "continue", "href": "https://juicer-dev.xyz/auth/login", "title": {"id": 1070009, "text": "Continue", "type": "info"}, "node_type": "a"}}], "action": "https://juicer-dev.xyz/auth/login", "method": "GET", "messages": [{"id": 1080002, "text": "You successfully verified your email address.", "type": "success"}]}	8712f917-6eee-4fc0-98b3-364070db5d96	1
\.


--
-- Data for Name: session_devices; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.session_devices (id, ip_address, user_agent, location, nid, session_id, created_at, updated_at) FROM stdin;
79fc8d02-26a2-4003-93a4-26c3f9ff9a81	172.18.0.1	Mozilla/5.0 (X11; Linux x86_64; rv:122.0) Gecko/20100101 Firefox/122.0		8712f917-6eee-4fc0-98b3-364070db5d96	1f8ac415-ab47-487c-aea3-bdb7b09bc52a	2024-02-09 16:00:48.814731	2024-02-09 16:00:48.814731
\.


--
-- Data for Name: session_token_exchanges; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.session_token_exchanges (id, nid, flow_id, session_id, init_code, return_to_code, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: sessions; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.sessions (id, issued_at, expires_at, authenticated_at, identity_id, created_at, updated_at, token, active, nid, logout_token, aal, authentication_methods) FROM stdin;
1f8ac415-ab47-487c-aea3-bdb7b09bc52a	2024-02-09 16:00:48.81095	2024-02-10 16:00:48.81095	2024-02-09 16:00:48.81095	2985b0ba-d796-427c-b1ae-bafb26cb25f0	2024-02-09 16:00:48.81219	2024-02-09 16:00:48.81219	ory_st_OKTxPVfZaIo4GyDf7cnXBtYEg3Zk5rah	f	8712f917-6eee-4fc0-98b3-364070db5d96	ory_lo_D60uATUnL0Mo4U72rUOkSeIGfdZgZvtu	aal1	[{"aal": "aal1", "method": "password", "completed_at": "2024-02-09T16:00:48.810947051Z"}]
\.


--
-- Name: continuity_containers continuity_containers_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.continuity_containers
    ADD CONSTRAINT continuity_containers_pkey PRIMARY KEY (id);


--
-- Name: courier_message_dispatches courier_message_dispatches_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.courier_message_dispatches
    ADD CONSTRAINT courier_message_dispatches_pkey PRIMARY KEY (id);


--
-- Name: courier_messages courier_messages_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.courier_messages
    ADD CONSTRAINT courier_messages_pkey PRIMARY KEY (id);


--
-- Name: identities identities_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identities
    ADD CONSTRAINT identities_pkey PRIMARY KEY (id);


--
-- Name: identity_credential_identifiers identity_credential_identifiers_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_credential_identifiers
    ADD CONSTRAINT identity_credential_identifiers_pkey PRIMARY KEY (id);


--
-- Name: identity_credential_types identity_credential_types_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_credential_types
    ADD CONSTRAINT identity_credential_types_pkey PRIMARY KEY (id);


--
-- Name: identity_credentials identity_credentials_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_credentials
    ADD CONSTRAINT identity_credentials_pkey PRIMARY KEY (id);


--
-- Name: identity_recovery_addresses identity_recovery_addresses_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_recovery_addresses
    ADD CONSTRAINT identity_recovery_addresses_pkey PRIMARY KEY (id);


--
-- Name: identity_recovery_codes identity_recovery_codes_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_recovery_codes
    ADD CONSTRAINT identity_recovery_codes_pkey PRIMARY KEY (id);


--
-- Name: identity_recovery_tokens identity_recovery_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_recovery_tokens
    ADD CONSTRAINT identity_recovery_tokens_pkey PRIMARY KEY (id);


--
-- Name: identity_verifiable_addresses identity_verifiable_addresses_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_verifiable_addresses
    ADD CONSTRAINT identity_verifiable_addresses_pkey PRIMARY KEY (id);


--
-- Name: identity_verification_codes identity_verification_codes_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_verification_codes
    ADD CONSTRAINT identity_verification_codes_pkey PRIMARY KEY (id);


--
-- Name: identity_verification_tokens identity_verification_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_verification_tokens
    ADD CONSTRAINT identity_verification_tokens_pkey PRIMARY KEY (id);


--
-- Name: keto_relation_tuples keto_relation_tuples_uuid_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.keto_relation_tuples
    ADD CONSTRAINT keto_relation_tuples_uuid_pkey PRIMARY KEY (shard_id, nid);


--
-- Name: keto_uuid_mappings keto_uuid_mappings_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.keto_uuid_mappings
    ADD CONSTRAINT keto_uuid_mappings_pkey PRIMARY KEY (id);


--
-- Name: networks networks_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.networks
    ADD CONSTRAINT networks_pkey PRIMARY KEY (id);


--
-- Name: selfservice_errors selfservice_errors_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_errors
    ADD CONSTRAINT selfservice_errors_pkey PRIMARY KEY (id);


--
-- Name: selfservice_login_flows selfservice_login_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_login_flows
    ADD CONSTRAINT selfservice_login_requests_pkey PRIMARY KEY (id);


--
-- Name: selfservice_settings_flows selfservice_profile_management_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_settings_flows
    ADD CONSTRAINT selfservice_profile_management_requests_pkey PRIMARY KEY (id);


--
-- Name: selfservice_recovery_flows selfservice_recovery_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_recovery_flows
    ADD CONSTRAINT selfservice_recovery_requests_pkey PRIMARY KEY (id);


--
-- Name: selfservice_registration_flows selfservice_registration_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_registration_flows
    ADD CONSTRAINT selfservice_registration_requests_pkey PRIMARY KEY (id);


--
-- Name: selfservice_verification_flows selfservice_verification_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_verification_flows
    ADD CONSTRAINT selfservice_verification_requests_pkey PRIMARY KEY (id);


--
-- Name: session_devices session_devices_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.session_devices
    ADD CONSTRAINT session_devices_pkey PRIMARY KEY (id);


--
-- Name: session_token_exchanges session_token_exchanges_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.session_token_exchanges
    ADD CONSTRAINT session_token_exchanges_pkey PRIMARY KEY (id);


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id);


--
-- Name: session_devices unique_session_device; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.session_devices
    ADD CONSTRAINT unique_session_device UNIQUE (nid, session_id, ip_address, user_agent);


--
-- Name: continuity_containers_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX continuity_containers_id_nid_idx ON public.continuity_containers USING btree (id, nid);


--
-- Name: continuity_containers_identity_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX continuity_containers_identity_id_nid_idx ON public.continuity_containers USING btree (identity_id, nid);


--
-- Name: continuity_containers_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX continuity_containers_nid_id_idx ON public.continuity_containers USING btree (nid, id);


--
-- Name: courier_message_dispatches_id_message_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX courier_message_dispatches_id_message_id_nid_idx ON public.courier_message_dispatches USING btree (id, message_id, nid);


--
-- Name: courier_messages_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX courier_messages_id_nid_idx ON public.courier_messages USING btree (id, nid);


--
-- Name: courier_messages_nid_created_at_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX courier_messages_nid_created_at_id_idx ON public.courier_messages USING btree (nid, created_at DESC);


--
-- Name: courier_messages_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX courier_messages_nid_id_idx ON public.courier_messages USING btree (nid, id);


--
-- Name: courier_messages_nid_recipient_created_at_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX courier_messages_nid_recipient_created_at_id_idx ON public.courier_messages USING btree (nid, recipient, created_at DESC);


--
-- Name: courier_messages_nid_status_created_at_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX courier_messages_nid_status_created_at_id_idx ON public.courier_messages USING btree (nid, status, created_at DESC);


--
-- Name: courier_messages_status_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX courier_messages_status_idx ON public.courier_messages USING btree (status);


--
-- Name: identities_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identities_id_nid_idx ON public.identities USING btree (id, nid);


--
-- Name: identities_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identities_nid_id_idx ON public.identities USING btree (nid, id);


--
-- Name: identity_credential_identifiers_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_credential_identifiers_id_nid_idx ON public.identity_credential_identifiers USING btree (id, nid);


--
-- Name: identity_credential_identifiers_identifier_nid_type_uq_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX identity_credential_identifiers_identifier_nid_type_uq_idx ON public.identity_credential_identifiers USING btree (nid, identity_credential_type_id, identifier);


--
-- Name: identity_credential_identifiers_nid_i_ici_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_credential_identifiers_nid_i_ici_idx ON public.identity_credential_identifiers USING btree (nid, identifier, identity_credential_id);


--
-- Name: identity_credential_identifiers_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_credential_identifiers_nid_id_idx ON public.identity_credential_identifiers USING btree (nid, id);


--
-- Name: identity_credential_identifiers_nid_identity_credential_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_credential_identifiers_nid_identity_credential_id_idx ON public.identity_credential_identifiers USING btree (identity_credential_id, nid);


--
-- Name: identity_credential_types_name_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX identity_credential_types_name_idx ON public.identity_credential_types USING btree (name);


--
-- Name: identity_credentials_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_credentials_id_nid_idx ON public.identity_credentials USING btree (id, nid);


--
-- Name: identity_credentials_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_credentials_nid_id_idx ON public.identity_credentials USING btree (nid, id);


--
-- Name: identity_credentials_nid_identity_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_credentials_nid_identity_id_idx ON public.identity_credentials USING btree (identity_id, nid);


--
-- Name: identity_recovery_addresses_code_uq_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX identity_recovery_addresses_code_uq_idx ON public.identity_recovery_tokens USING btree (token);


--
-- Name: identity_recovery_addresses_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_addresses_id_nid_idx ON public.identity_recovery_addresses USING btree (id, nid);


--
-- Name: identity_recovery_addresses_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_addresses_nid_id_idx ON public.identity_recovery_addresses USING btree (nid, id);


--
-- Name: identity_recovery_addresses_nid_identity_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_addresses_nid_identity_id_idx ON public.identity_recovery_addresses USING btree (identity_id, nid);


--
-- Name: identity_recovery_addresses_status_via_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_addresses_status_via_idx ON public.identity_recovery_addresses USING btree (nid, via, value);


--
-- Name: identity_recovery_addresses_status_via_uq_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX identity_recovery_addresses_status_via_uq_idx ON public.identity_recovery_addresses USING btree (nid, via, value);


--
-- Name: identity_recovery_codes_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_codes_id_nid_idx ON public.identity_recovery_codes USING btree (id, nid);


--
-- Name: identity_recovery_codes_identity_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_codes_identity_id_nid_idx ON public.identity_recovery_codes USING btree (identity_id, nid);


--
-- Name: identity_recovery_codes_identity_recovery_address_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_codes_identity_recovery_address_id_nid_idx ON public.identity_recovery_codes USING btree (identity_recovery_address_id, nid);


--
-- Name: identity_recovery_codes_nid_flow_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_codes_nid_flow_id_idx ON public.identity_recovery_codes USING btree (nid, selfservice_recovery_flow_id);


--
-- Name: identity_recovery_tokens_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_tokens_id_nid_idx ON public.identity_recovery_tokens USING btree (id, nid);


--
-- Name: identity_recovery_tokens_identity_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_tokens_identity_id_nid_idx ON public.identity_recovery_tokens USING btree (identity_id, nid);


--
-- Name: identity_recovery_tokens_identity_recovery_address_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_tokens_identity_recovery_address_id_idx ON public.identity_recovery_tokens USING btree (identity_recovery_address_id);


--
-- Name: identity_recovery_tokens_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_tokens_nid_id_idx ON public.identity_recovery_tokens USING btree (nid, id);


--
-- Name: identity_recovery_tokens_selfservice_recovery_flow_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_tokens_selfservice_recovery_flow_id_idx ON public.identity_recovery_tokens USING btree (selfservice_recovery_flow_id);


--
-- Name: identity_recovery_tokens_token_nid_used_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_tokens_token_nid_used_idx ON public.identity_recovery_tokens USING btree (nid, token, used);


--
-- Name: identity_verifiable_addresses_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_verifiable_addresses_id_nid_idx ON public.identity_verifiable_addresses USING btree (id, nid);


--
-- Name: identity_verifiable_addresses_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_verifiable_addresses_nid_id_idx ON public.identity_verifiable_addresses USING btree (nid, id);


--
-- Name: identity_verifiable_addresses_nid_identity_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_verifiable_addresses_nid_identity_id_idx ON public.identity_verifiable_addresses USING btree (identity_id, nid);


--
-- Name: identity_verifiable_addresses_status_via_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_verifiable_addresses_status_via_idx ON public.identity_verifiable_addresses USING btree (nid, via, value);


--
-- Name: identity_verifiable_addresses_status_via_uq_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX identity_verifiable_addresses_status_via_uq_idx ON public.identity_verifiable_addresses USING btree (nid, via, value);


--
-- Name: identity_verification_codes_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_verification_codes_id_nid_idx ON public.identity_verification_codes USING btree (id, nid);


--
-- Name: identity_verification_codes_nid_flow_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_verification_codes_nid_flow_id_idx ON public.identity_verification_codes USING btree (nid, selfservice_verification_flow_id);


--
-- Name: identity_verification_codes_verifiable_address_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_verification_codes_verifiable_address_nid_idx ON public.identity_verification_codes USING btree (identity_verifiable_address_id, nid);


--
-- Name: identity_verification_tokens_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_verification_tokens_id_nid_idx ON public.identity_verification_tokens USING btree (id, nid);


--
-- Name: identity_verification_tokens_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_verification_tokens_nid_id_idx ON public.identity_verification_tokens USING btree (nid, id);


--
-- Name: identity_verification_tokens_token_nid_used_flow_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_verification_tokens_token_nid_used_flow_id_idx ON public.identity_verification_tokens USING btree (nid, token, used, selfservice_verification_flow_id);


--
-- Name: identity_verification_tokens_token_uq_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX identity_verification_tokens_token_uq_idx ON public.identity_verification_tokens USING btree (token);


--
-- Name: identity_verification_tokens_verifiable_address_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_verification_tokens_verifiable_address_id_idx ON public.identity_verification_tokens USING btree (identity_verifiable_address_id);


--
-- Name: identity_verification_tokens_verification_flow_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_verification_tokens_verification_flow_id_idx ON public.identity_verification_tokens USING btree (selfservice_verification_flow_id);


--
-- Name: keto_relation_tuples_uuid_full_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX keto_relation_tuples_uuid_full_idx ON public.keto_relation_tuples USING btree (nid, namespace, object, relation, subject_id, subject_set_namespace, subject_set_object, subject_set_relation, commit_time);


--
-- Name: keto_relation_tuples_uuid_reverse_subject_ids_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX keto_relation_tuples_uuid_reverse_subject_ids_idx ON public.keto_relation_tuples USING btree (nid, subject_id, relation, namespace) WHERE ((subject_set_namespace IS NULL) AND (subject_set_object IS NULL) AND (subject_set_relation IS NULL));


--
-- Name: keto_relation_tuples_uuid_reverse_subject_sets_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX keto_relation_tuples_uuid_reverse_subject_sets_idx ON public.keto_relation_tuples USING btree (nid, subject_set_namespace, subject_set_object, subject_set_relation, relation, namespace) WHERE (subject_id IS NULL);


--
-- Name: keto_relation_tuples_uuid_subject_ids_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX keto_relation_tuples_uuid_subject_ids_idx ON public.keto_relation_tuples USING btree (nid, namespace, object, relation, subject_id) WHERE ((subject_set_namespace IS NULL) AND (subject_set_object IS NULL) AND (subject_set_relation IS NULL));


--
-- Name: keto_relation_tuples_uuid_subject_sets_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX keto_relation_tuples_uuid_subject_sets_idx ON public.keto_relation_tuples USING btree (nid, namespace, object, relation, subject_set_namespace, subject_set_object, subject_set_relation) WHERE (subject_id IS NULL);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- Name: schema_migration_version_self_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX schema_migration_version_self_idx ON public.schema_migration USING btree (version_self);


--
-- Name: selfservice_errors_errors_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX selfservice_errors_errors_nid_id_idx ON public.selfservice_errors USING btree (nid, id);


--
-- Name: selfservice_login_flows_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX selfservice_login_flows_id_nid_idx ON public.selfservice_login_flows USING btree (id, nid);


--
-- Name: selfservice_login_flows_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX selfservice_login_flows_nid_id_idx ON public.selfservice_login_flows USING btree (nid, id);


--
-- Name: selfservice_recovery_flows_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX selfservice_recovery_flows_id_nid_idx ON public.selfservice_recovery_flows USING btree (id, nid);


--
-- Name: selfservice_recovery_flows_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX selfservice_recovery_flows_nid_id_idx ON public.selfservice_recovery_flows USING btree (nid, id);


--
-- Name: selfservice_recovery_flows_recovered_identity_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX selfservice_recovery_flows_recovered_identity_id_nid_idx ON public.selfservice_recovery_flows USING btree (recovered_identity_id, nid);


--
-- Name: selfservice_registration_flows_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX selfservice_registration_flows_id_nid_idx ON public.selfservice_registration_flows USING btree (id, nid);


--
-- Name: selfservice_registration_flows_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX selfservice_registration_flows_nid_id_idx ON public.selfservice_registration_flows USING btree (nid, id);


--
-- Name: selfservice_settings_flows_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX selfservice_settings_flows_id_nid_idx ON public.selfservice_settings_flows USING btree (id, nid);


--
-- Name: selfservice_settings_flows_identity_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX selfservice_settings_flows_identity_id_nid_idx ON public.selfservice_settings_flows USING btree (identity_id, nid);


--
-- Name: selfservice_settings_flows_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX selfservice_settings_flows_nid_id_idx ON public.selfservice_settings_flows USING btree (nid, id);


--
-- Name: selfservice_verification_flows_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX selfservice_verification_flows_id_nid_idx ON public.selfservice_verification_flows USING btree (id, nid);


--
-- Name: selfservice_verification_flows_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX selfservice_verification_flows_nid_id_idx ON public.selfservice_verification_flows USING btree (nid, id);


--
-- Name: session_devices_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX session_devices_id_nid_idx ON public.session_devices USING btree (id, nid);


--
-- Name: session_devices_session_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX session_devices_session_id_nid_idx ON public.session_devices USING btree (session_id, nid);


--
-- Name: session_token_exchanges_nid_code_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX session_token_exchanges_nid_code_idx ON public.session_token_exchanges USING btree (init_code, nid);


--
-- Name: session_token_exchanges_nid_flow_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX session_token_exchanges_nid_flow_id_idx ON public.session_token_exchanges USING btree (flow_id, nid);


--
-- Name: sessions_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX sessions_id_nid_idx ON public.sessions USING btree (id, nid);


--
-- Name: sessions_identity_id_nid_sorted_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX sessions_identity_id_nid_sorted_idx ON public.sessions USING btree (identity_id, nid, authenticated_at DESC);


--
-- Name: sessions_logout_token_uq_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX sessions_logout_token_uq_idx ON public.sessions USING btree (logout_token);


--
-- Name: sessions_nid_id_identity_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX sessions_nid_id_identity_id_idx ON public.sessions USING btree (nid, identity_id, id);


--
-- Name: sessions_token_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX sessions_token_nid_idx ON public.sessions USING btree (nid, token);


--
-- Name: sessions_token_uq_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX sessions_token_uq_idx ON public.sessions USING btree (token);


--
-- Name: continuity_containers continuity_containers_identity_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.continuity_containers
    ADD CONSTRAINT continuity_containers_identity_id_fkey FOREIGN KEY (identity_id) REFERENCES public.identities(id) ON DELETE CASCADE;


--
-- Name: continuity_containers continuity_containers_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.continuity_containers
    ADD CONSTRAINT continuity_containers_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: courier_message_dispatches courier_message_dispatches_message_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.courier_message_dispatches
    ADD CONSTRAINT courier_message_dispatches_message_id_fk FOREIGN KEY (message_id) REFERENCES public.courier_messages(id) ON DELETE CASCADE;


--
-- Name: courier_message_dispatches courier_message_dispatches_nid_fk; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.courier_message_dispatches
    ADD CONSTRAINT courier_message_dispatches_nid_fk FOREIGN KEY (nid) REFERENCES public.networks(id) ON DELETE CASCADE;


--
-- Name: courier_messages courier_messages_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.courier_messages
    ADD CONSTRAINT courier_messages_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identities identities_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identities
    ADD CONSTRAINT identities_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identity_credential_identifiers identity_credential_identifiers_identity_credential_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_credential_identifiers
    ADD CONSTRAINT identity_credential_identifiers_identity_credential_id_fkey FOREIGN KEY (identity_credential_id) REFERENCES public.identity_credentials(id) ON DELETE CASCADE;


--
-- Name: identity_credential_identifiers identity_credential_identifiers_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_credential_identifiers
    ADD CONSTRAINT identity_credential_identifiers_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identity_credential_identifiers identity_credential_identifiers_type_id_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_credential_identifiers
    ADD CONSTRAINT identity_credential_identifiers_type_id_fk_idx FOREIGN KEY (identity_credential_type_id) REFERENCES public.identity_credential_types(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identity_credentials identity_credentials_identity_credential_type_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_credentials
    ADD CONSTRAINT identity_credentials_identity_credential_type_id_fkey FOREIGN KEY (identity_credential_type_id) REFERENCES public.identity_credential_types(id) ON DELETE CASCADE;


--
-- Name: identity_credentials identity_credentials_identity_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_credentials
    ADD CONSTRAINT identity_credentials_identity_id_fkey FOREIGN KEY (identity_id) REFERENCES public.identities(id) ON DELETE CASCADE;


--
-- Name: identity_credentials identity_credentials_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_credentials
    ADD CONSTRAINT identity_credentials_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identity_recovery_addresses identity_recovery_addresses_identity_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_recovery_addresses
    ADD CONSTRAINT identity_recovery_addresses_identity_id_fkey FOREIGN KEY (identity_id) REFERENCES public.identities(id) ON DELETE CASCADE;


--
-- Name: identity_recovery_addresses identity_recovery_addresses_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_recovery_addresses
    ADD CONSTRAINT identity_recovery_addresses_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identity_recovery_codes identity_recovery_codes_identity_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_recovery_codes
    ADD CONSTRAINT identity_recovery_codes_identity_id_fk FOREIGN KEY (identity_id) REFERENCES public.identities(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identity_recovery_codes identity_recovery_codes_identity_recovery_addresses_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_recovery_codes
    ADD CONSTRAINT identity_recovery_codes_identity_recovery_addresses_id_fk FOREIGN KEY (identity_recovery_address_id) REFERENCES public.identity_recovery_addresses(id) ON DELETE CASCADE;


--
-- Name: identity_recovery_codes identity_recovery_codes_networks_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_recovery_codes
    ADD CONSTRAINT identity_recovery_codes_networks_id_fk FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identity_recovery_codes identity_recovery_codes_selfservice_recovery_flows_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_recovery_codes
    ADD CONSTRAINT identity_recovery_codes_selfservice_recovery_flows_id_fk FOREIGN KEY (selfservice_recovery_flow_id) REFERENCES public.selfservice_recovery_flows(id) ON DELETE CASCADE;


--
-- Name: identity_recovery_tokens identity_recovery_tokens_identity_id_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_recovery_tokens
    ADD CONSTRAINT identity_recovery_tokens_identity_id_fk_idx FOREIGN KEY (identity_id) REFERENCES public.identities(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identity_recovery_tokens identity_recovery_tokens_identity_recovery_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_recovery_tokens
    ADD CONSTRAINT identity_recovery_tokens_identity_recovery_address_id_fkey FOREIGN KEY (identity_recovery_address_id) REFERENCES public.identity_recovery_addresses(id) ON DELETE CASCADE;


--
-- Name: identity_recovery_tokens identity_recovery_tokens_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_recovery_tokens
    ADD CONSTRAINT identity_recovery_tokens_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identity_recovery_tokens identity_recovery_tokens_selfservice_recovery_request_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_recovery_tokens
    ADD CONSTRAINT identity_recovery_tokens_selfservice_recovery_request_id_fkey FOREIGN KEY (selfservice_recovery_flow_id) REFERENCES public.selfservice_recovery_flows(id) ON DELETE CASCADE;


--
-- Name: identity_verifiable_addresses identity_verifiable_addresses_identity_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_verifiable_addresses
    ADD CONSTRAINT identity_verifiable_addresses_identity_id_fkey FOREIGN KEY (identity_id) REFERENCES public.identities(id) ON DELETE CASCADE;


--
-- Name: identity_verifiable_addresses identity_verifiable_addresses_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_verifiable_addresses
    ADD CONSTRAINT identity_verifiable_addresses_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identity_verification_codes identity_verification_codes_identity_verifiable_addresses_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_verification_codes
    ADD CONSTRAINT identity_verification_codes_identity_verifiable_addresses_id_fk FOREIGN KEY (identity_verifiable_address_id) REFERENCES public.identity_verifiable_addresses(id) ON DELETE CASCADE;


--
-- Name: identity_verification_codes identity_verification_codes_networks_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_verification_codes
    ADD CONSTRAINT identity_verification_codes_networks_id_fk FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identity_verification_codes identity_verification_codes_selfservice_verification_flows_id_f; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_verification_codes
    ADD CONSTRAINT identity_verification_codes_selfservice_verification_flows_id_f FOREIGN KEY (selfservice_verification_flow_id) REFERENCES public.selfservice_verification_flows(id) ON DELETE CASCADE;


--
-- Name: identity_verification_tokens identity_verification_tokens_identity_verifiable_address_i_fkey; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_verification_tokens
    ADD CONSTRAINT identity_verification_tokens_identity_verifiable_address_i_fkey FOREIGN KEY (identity_verifiable_address_id) REFERENCES public.identity_verifiable_addresses(id) ON DELETE CASCADE;


--
-- Name: identity_verification_tokens identity_verification_tokens_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_verification_tokens
    ADD CONSTRAINT identity_verification_tokens_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identity_verification_tokens identity_verification_tokens_selfservice_verification_flow_fkey; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_verification_tokens
    ADD CONSTRAINT identity_verification_tokens_selfservice_verification_flow_fkey FOREIGN KEY (selfservice_verification_flow_id) REFERENCES public.selfservice_verification_flows(id) ON DELETE CASCADE;


--
-- Name: keto_relation_tuples keto_relation_tuples_nid_fk; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.keto_relation_tuples
    ADD CONSTRAINT keto_relation_tuples_nid_fk FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: selfservice_errors selfservice_errors_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_errors
    ADD CONSTRAINT selfservice_errors_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: selfservice_login_flows selfservice_login_flows_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_login_flows
    ADD CONSTRAINT selfservice_login_flows_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: selfservice_settings_flows selfservice_profile_management_requests_identity_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_settings_flows
    ADD CONSTRAINT selfservice_profile_management_requests_identity_id_fkey FOREIGN KEY (identity_id) REFERENCES public.identities(id) ON DELETE CASCADE;


--
-- Name: selfservice_recovery_flows selfservice_recovery_flows_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_recovery_flows
    ADD CONSTRAINT selfservice_recovery_flows_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: selfservice_recovery_flows selfservice_recovery_requests_recovered_identity_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_recovery_flows
    ADD CONSTRAINT selfservice_recovery_requests_recovered_identity_id_fkey FOREIGN KEY (recovered_identity_id) REFERENCES public.identities(id) ON DELETE CASCADE;


--
-- Name: selfservice_registration_flows selfservice_registration_flows_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_registration_flows
    ADD CONSTRAINT selfservice_registration_flows_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: selfservice_settings_flows selfservice_settings_flows_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_settings_flows
    ADD CONSTRAINT selfservice_settings_flows_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: selfservice_verification_flows selfservice_verification_flows_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_verification_flows
    ADD CONSTRAINT selfservice_verification_flows_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: session_devices session_metadata_nid_fk; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.session_devices
    ADD CONSTRAINT session_metadata_nid_fk FOREIGN KEY (nid) REFERENCES public.networks(id) ON DELETE CASCADE;


--
-- Name: session_devices session_metadata_sessions_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.session_devices
    ADD CONSTRAINT session_metadata_sessions_id_fk FOREIGN KEY (session_id) REFERENCES public.sessions(id) ON DELETE CASCADE;


--
-- Name: sessions sessions_identity_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_identity_id_fkey FOREIGN KEY (identity_id) REFERENCES public.identities(id) ON DELETE CASCADE;


--
-- Name: sessions sessions_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

--
-- PostgreSQL database cluster dump complete
--

