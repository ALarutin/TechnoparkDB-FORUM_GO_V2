CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;

------------------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------
-- TABLES --------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------

-- table person
CREATE TABLE public.person
(
    id       SERIAL NOT NULL,
    nickname citext NOT NULL,
    email    citext NOT NULL,
    fullname text   NOT NULL,
    about    text   NOT NULL
);

CREATE UNIQUE INDEX person_email_ui
    ON public.person (email);

ALTER TABLE public.person
    ADD CONSTRAINT person_pk PRIMARY KEY (nickname);

-- table forum
CREATE TABLE public.forum
(
    id      SERIAL          NOT NULL,
    slug    citext          NOT NULL,
    author  citext          NOT NULL,
    title   text DEFAULT '' NOT NULL,
    posts   INT  DEFAULT 0  NOT NULL,
    threads INT  DEFAULT 0  NOT NULL
);

ALTER TABLE public.forum
    ADD CONSTRAINT forum_pk PRIMARY KEY (slug);

ALTER TABLE ONLY public.forum
    ADD CONSTRAINT "forum_user_fk" FOREIGN KEY (author) REFERENCES public.person (nickname);

-- table forum_users
CREATE TABLE public.forum_users
(
    forum_slug    citext NOT NULL,
    user_nickname citext NOT NULL
);

ALTER TABLE ONLY public.forum_users
    ADD CONSTRAINT "forum_users_forum_slug_fk" FOREIGN KEY (forum_slug) REFERENCES public.forum (slug);

ALTER TABLE ONLY public.forum_users
    ADD CONSTRAINT "forum_users_user_nickname_fk" FOREIGN KEY (user_nickname) REFERENCES public.person (nickname);

-- table thread
CREATE TABLE public.thread
(
    id      SERIAL                   NOT NULL,
    slug    citext,
    author  citext                   NOT NULL,
    forum   citext                   NOT NULL,
    title   text DEFAULT ''          NOT NULL,
    message text DEFAULT ''          NOT NULL,
    votes   INT  DEFAULT 0           NOT NULL,
    created TIMESTAMP WITH TIME ZONE NOT NULL
);

ALTER TABLE thread
    ADD CONSTRAINT thread_pk PRIMARY KEY (id);

CREATE UNIQUE INDEX thread_slug_ui
    ON public.thread (slug);

ALTER TABLE ONLY public.thread
    ADD CONSTRAINT "thread_author_fk" FOREIGN KEY (author) REFERENCES public.person (nickname);

ALTER TABLE ONLY public.thread
    ADD CONSTRAINT "thread_forum_fk" FOREIGN KEY (forum) REFERENCES public.forum (slug);

-- table post
CREATE TABLE public.post
(
    id        SERIAL                                             NOT NULL,
    author    citext                                             NOT NULL,
    thread    INT                                                NOT NULL,
    forum     citext                                             NOT NULL,
    message   text                     DEFAULT ''                NOT NULL,
    is_edited BOOLEAN                  DEFAULT FALSE             NOT NULL,
    parent    INT                                                NOT NULL,
    created   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    post_path INT[]                    DEFAULT '{}'::INT[]
);

ALTER TABLE public.post
    ADD CONSTRAINT post_pk PRIMARY KEY (id);

ALTER TABLE ONLY public.post
    ADD CONSTRAINT "post_author_fk" FOREIGN KEY (author) REFERENCES public.person (nickname);

ALTER TABLE ONLY public.post
    ADD CONSTRAINT "post_thread_fk" FOREIGN KEY (thread) REFERENCES public.thread (id);

ALTER TABLE ONLY public.post
    ADD CONSTRAINT "post_forum_fk" FOREIGN KEY (forum) REFERENCES public.forum (slug);

ALTER TABLE ONLY public.post
    ADD CONSTRAINT "post_parent_fk" FOREIGN KEY (parent) REFERENCES public.post (id);

-- table vote
CREATE TABLE public.vote
(
    thread_id     INT    NOT NULL,
    user_nickname citext NOT NULL,
    voice         INT    NOT NULL
);

ALTER TABLE public.vote
    ADD CONSTRAINT vote_pk PRIMARY KEY (thread_id, user_nickname);

ALTER TABLE ONLY public.vote
    ADD CONSTRAINT "vote_thread_slug_fk" FOREIGN KEY (thread_id) REFERENCES public.thread (id);

ALTER TABLE ONLY public.vote
    ADD CONSTRAINT "vote_user_nickname_fk" FOREIGN KEY (user_nickname) REFERENCES public.person (nickname);

------------------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------
-- TRIGGERS ------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION update_threads_quantity()
    RETURNS trigger AS
$BODY$
BEGIN
    UPDATE public."forum"
    SET threads = threads + 1
    WHERE "slug" = NEW."forum";
    RETURN NULL;
END;
$BODY$
    LANGUAGE plpgsql;

CREATE TRIGGER update_forum_threads
    AFTER INSERT
    ON public.thread
    FOR EACH ROW
EXECUTE PROCEDURE update_threads_quantity();

CREATE OR REPLACE FUNCTION update_forum_users_on_thread()
    RETURNS trigger AS
$BODY$
BEGIN
    INSERT INTO public."forum_users"(forum_slug, user_nickname)
    VALUES (NEW."forum", NEW."author");
    RETURN NULL;
END;
$BODY$
    LANGUAGE plpgsql;

CREATE TRIGGER update_forum_users_on_thread
    AFTER INSERT
    ON public.thread
    FOR EACH ROW
EXECUTE PROCEDURE update_forum_users_on_thread();

CREATE OR REPLACE FUNCTION update_forum_users_on_post()
    RETURNS trigger AS
$BODY$
BEGIN
    INSERT INTO public."forum_users"(forum_slug, user_nickname)
    VALUES (NEW."forum", NEW."author");
    RETURN NULL;
END;
$BODY$
    LANGUAGE plpgsql;

CREATE TRIGGER update_forum_users_on_thread
    AFTER INSERT
    ON public.post
    FOR EACH ROW
EXECUTE PROCEDURE update_forum_users_on_thread();

CREATE OR REPLACE FUNCTION update_posts_quantity()
    RETURNS trigger AS
$BODY$
BEGIN
    UPDATE public."forum"
    SET posts = posts + 1
    WHERE "slug" = NEW."forum";
    RETURN NULL;
END;
$BODY$
    LANGUAGE plpgsql;

CREATE TRIGGER update_forum_posts
    AFTER INSERT
    ON public.post
    FOR EACH ROW
EXECUTE PROCEDURE update_posts_quantity();

CREATE OR REPLACE FUNCTION update_post_path()
    RETURNS trigger AS
$BODY$
DECLARE
    arg_post_path INT[];
BEGIN
    SELECT post_path
    INTO arg_post_path
    FROM public.post
    WHERE id = NEW.parent;
    arg_post_path = arg_post_path || ARRAY [New.id];
    UPDATE public.post
    SET post_path = arg_post_path
    WHERE id = NEW.id;
    RETURN NULL;
END;
$BODY$
    LANGUAGE plpgsql;

CREATE TRIGGER update_post_path
    AFTER INSERT
    ON public.post
    FOR EACH ROW
EXECUTE PROCEDURE update_post_path();

CREATE OR REPLACE FUNCTION insert_votes()
    RETURNS trigger AS
$BODY$
BEGIN
    UPDATE public."thread"
    SET votes = votes + New."voice"
    WHERE "id" = NEW."thread_id";
    RETURN NULL;
END;
$BODY$
    LANGUAGE plpgsql;

CREATE TRIGGER insert_thread_votes
    AFTER INSERT
    ON public.vote
    FOR EACH ROW
EXECUTE PROCEDURE insert_votes();

CREATE OR REPLACE FUNCTION update_votes()
    RETURNS trigger AS
$BODY$
BEGIN
    IF (OLD.voice != NEW.voice) THEN
        UPDATE public."thread"
        SET votes = votes + 2 * New.voice
        WHERE "id" = NEW.thread_id;
    END IF;
    RETURN NULL;
END;
$BODY$
    LANGUAGE plpgsql;

CREATE TRIGGER update_thread_votes
    AFTER UPDATE
    ON public.vote
    FOR EACH ROW
EXECUTE PROCEDURE update_votes();

------------------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------
-- TYPES ---------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------

CREATE TYPE public.type_person AS
    (
    is_new BOOLEAN,
    id BIGINT,
    nickname citext,
    email citext,
    fullname text,
    about text
    );

CREATE TYPE public.type_forum AS
    (
    is_new BOOLEAN,
    id BIGINT,
    slug citext,
    author citext,
    title text,
    posts INT,
    threads INT
    );

CREATE TYPE public.type_thread AS
    (
    is_new BOOLEAN,
    id BIGINT,
    slug citext,
    author citext,
    forum citext,
    title text,
    message text,
    votes INT,
    created TIMESTAMP WITH TIME ZONE
    );

CREATE TYPE public.type_post AS
    (
    id BIGINT,
    author citext,
    thread INT,
    forum citext,
    message text,
    is_edited BOOLEAN,
    parent INT,
    created TIMESTAMP WITH TIME ZONE,
    post_path INT[]
    );

CREATE TYPE public.type_database AS
    (
    forum INT,
    post INT,
    thread INT,
    person INT
    );

------------------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------
-- FUNCTIONS -----------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION func_add_admin()
    RETURNS VOID AS
$BODY$
BEGIN
    INSERT INTO public."person" (id, email, about, fullname, nickname)
    VALUES (0, 'admin@admin.com', 'something', 'admin', 'admin');
    INSERT INTO public."forum" (id, author, slug)
    VALUES (0, 'admin', 'admin');
    INSERT INTO public."thread" (id, author, forum, slug, created)
    VALUES (0, 'admin', 'admin', 'admin', '0001-01-01 00:00:00.000000 +00:00');
    INSERT INTO public."post" (id, author, thread, forum, parent, created)
    VALUES (0, 'admin', '0', 'admin', 0, '0001-01-01 00:00:00.000000 +00:00');
END;
$BODY$
    LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION func_clear_database()
    RETURNS VOID AS
$BODY$
BEGIN
    TRUNCATE TABLE public.forum, public.forum_users, public.person, public.post, public.thread, public.vote
        RESTART IDENTITY;
    PERFORM func_add_admin();
END;
$BODY$
    LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION func_get_database()
    RETURNS public.type_database
AS
$BODY$
DECLARE
    result public.type_database;
BEGIN
    SELECT count(*) INTO result.person FROM public.person;
    SELECT count(*) INTO result.forum FROM public.forum;
    SELECT count(*) INTO result.thread FROM public.thread;
    SELECT count(*) INTO result.post FROM public.post;
    RETURN result;
END;
$BODY$
    LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION func_create_user(arg_nickname citext, arg_email citext, arg_fullname text, arg_about text)
    RETURNS SETOF public.type_person
AS
$BODY$
DECLARE
    result public.type_person;
    rec    RECORD;
BEGIN
    INSERT INTO person (nickname, email, fullname, about)
    VALUES (arg_nickname, arg_email, arg_fullname, arg_about) RETURNING *
        INTO result.id, result.nickname, result.email, result.fullname, result.about;
    result.is_new := true;
    RETURN next result;
EXCEPTION
    WHEN unique_violation THEN
        FOR rec IN SELECT *
                   FROM public.person
                   WHERE nickname = arg_nickname
                      OR email = arg_email
            LOOP
                result.id := rec.id;
                result.nickname := rec.nickname;
                result.fullname := rec.fullname;
                result.about := rec.about;
                result.email := rec.email;
                result.is_new := false;
                RETURN NEXT result;
            END LOOP;
END;
$BODY$
    LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION func_get_user(arg_nickname citext)
    RETURNS public.type_person
AS
$BODY$
DECLARE
    result public.type_person;
BEGIN
    SELECT *
    INTO result.id, result.nickname, result.email, result.fullname, result.about
    FROM public.person
    WHERE nickname = arg_nickname;
    result.is_new := FALSE;
    IF NOT FOUND THEN
        RAISE no_data_found;
    END IF;
    RETURN result;
END;
$BODY$
    LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION func_update_user(arg_nickname citext, arg_email citext, arg_fullname text, arg_about text)
    RETURNS public.type_person
AS
$BODY$
DECLARE
    result public.type_person;
BEGIN
    UPDATE public.person
    SET email    = CASE
                       WHEN arg_email != '' THEN arg_email
                       ELSE email END,
        fullname = CASE
                       WHEN arg_fullname != '' THEN arg_fullname
                       ELSE fullname END,
        about    = CASE
                       WHEN arg_about != '' THEN arg_about
                       ELSE about END
    WHERE nickname = arg_nickname RETURNING *
        INTO result.id, result.nickname, result.email, result.fullname, result.about;
    result.is_new := FALSE;
    IF NOT FOUND THEN
        RAISE no_data_found;
    END IF;
    RETURN result;
EXCEPTION
    WHEN unique_violation THEN
        RAISE unique_violation;
END;
$BODY$
    LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION func_create_forum(arg_author citext, arg_slug citext, arg_title text)
    RETURNS public.type_forum
AS
$BODY$
DECLARE
    result       public.type_forum;
    arg_nickname citext;
BEGIN
    SELECT nickname INTO arg_nickname FROM public.person WHERE nickname = arg_author;
    IF NOT FOUND THEN
        RAISE no_data_found;
    END IF;
    INSERT INTO public.forum (slug, author, title)
    VALUES (arg_slug, arg_nickname, arg_title) RETURNING *
        INTO result.id, result.slug, result.author, result.title, result.posts, result.threads;
    result.is_new := TRUE;
    RETURN result;
EXCEPTION
    WHEN unique_violation THEN
        BEGIN
            SELECT *
            INTO result.id, result.slug, result.author, result.title, result.posts, result.threads
            FROM public.forum f
            WHERE f.slug = arg_slug;
            result.is_new := FALSE;
            RETURN result;
        END;
    WHEN foreign_key_violation THEN
        RAISE no_data_found;
END;
$BODY$
    LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION func_create_thread(arg_author citext, arg_created TIMESTAMP WITH TIME ZONE, arg_forum citext,
                                              arg_message text, arg_slug citext, arg_title text)
    RETURNS public.type_thread
AS
$BODY$
DECLARE
    result         public.type_thread;
    arg_slug_forum citext;
BEGIN
    SELECT slug INTO arg_slug_forum FROM public.forum WHERE slug = arg_forum;
    IF NOT FOUND THEN
        RAISE no_data_found;
    END IF;
    INSERT INTO public.thread (slug, author, forum, title, message, created)
    VALUES (CASE WHEN arg_slug != '' THEN arg_slug ELSE NULL END, arg_author, arg_slug_forum, arg_title, arg_message,
            arg_created) RETURNING *
               INTO result.id, result.slug, result.author, result.forum, result.title, result.message, result.votes, result.created;
    result.is_new := TRUE;
    IF result.slug IS NULL
    THEN
        result.slug = '';
    END IF;
    RETURN result;
EXCEPTION
    WHEN unique_violation THEN
        BEGIN
            SELECT *
            INTO result.id, result.slug, result.author, result.forum, result.title, result.message, result.votes, result.created
            FROM public.thread t
            WHERE t.slug = arg_slug;
            result.is_new := FALSE;
            RETURN result;
        END;
    WHEN foreign_key_violation THEN
        RAISE no_data_found;
END;
$BODY$
    LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION func_get_forum(arg_slug citext)
    RETURNS public.type_forum
AS
$BODY$
DECLARE
    result public.type_forum;
BEGIN
    SELECT *
    INTO result.id, result.slug, result.author, result.title, result.posts, result.threads
    FROM public.forum
    WHERE slug = arg_slug;
    result.is_new := TRUE;
    IF NOT FOUND THEN
        RAISE no_data_found;
    END IF;
    RETURN result;
END;
$BODY$
    LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION func_get_threads(arg_slug citext, arg_since TIMESTAMP WITH TIME ZONE, arg_desc BOOLEAN,
                                            arg_limit INT)
    RETURNS SETOF public.type_thread
AS
$BODY$
DECLARE
    result     public.type_thread;
    forum_slug citext;
    rec        RECORD;
BEGIN
    SELECT slug
    INTO forum_slug
    FROM public.forum
    WHERE slug = arg_slug;
    IF NOT FOUND THEN
        RAISE no_data_found;
    END IF;
    FOR rec IN SELECT *
               FROM public.thread
               WHERE forum = forum_slug
                 AND CASE
                         WHEN arg_since = '0001-01-01 00:00:00.000000 +00:00' THEN TRUE
                         WHEN arg_desc THEN created <= arg_since
                         ELSE created >= arg_since END
               ORDER BY (CASE WHEN arg_desc THEN created END) DESC,
                        (CASE WHEN NOT arg_desc THEN created END) ASC
               LIMIT arg_limit
        LOOP
            result.is_new := false;
            result.id := rec.id;
            result.slug := rec.slug;
            result.author := rec.author;
            result.forum := rec.forum;
            result.title := rec.title;
            result.message := rec.message;
            result.votes := rec.votes;
            result.created := rec.created;
            IF result.slug IS NULL
            THEN
                result.slug = '';
            END IF;
            RETURN next result;
        END LOOP;
END;
$BODY$
    LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION func_get_users(arg_slug citext, arg_since citext, arg_desc BOOLEAN, arg_limit INT)
    RETURNS SETOF public.type_person
AS
$BODY$
DECLARE
    result public.type_person;
    rec    RECORD;
BEGIN
    PERFORM func_get_forum(arg_slug);
    FOR rec IN SELECT *
               FROM public.person
               WHERE nickname IN (SELECT user_nickname
                                  FROM public.forum_users
                                  WHERE forum_slug = arg_slug)
                 AND CASE
                         WHEN arg_since = '' THEN true
                         WHEN arg_desc THEN nickname < arg_since
                         ELSE nickname > arg_since END
               ORDER BY (CASE WHEN arg_desc THEN nickname END) DESC,
                        (CASE WHEN NOT arg_desc THEN nickname END) ASC
               LIMIT arg_limit
        LOOP
            result.is_new := false;
            result.id := rec.id;
            result.nickname := rec.nickname;
            result.email := rec.email;
            result.fullname := rec.fullname;
            result.about := rec.about;
            RETURN next result;
        END LOOP;
END;
$BODY$
    LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION func_get_thread_by_id(arg_id INT)
    RETURNS public.type_thread
AS
$BODY$
DECLARE
    result public.type_thread;
BEGIN
    SELECT *
    INTO result.id, result.slug, result.author, result.forum,
        result.title, result.message, result.votes, result.created
    FROM public.thread
    WHERE id = arg_id;
    result.is_new := FALSE;
    IF NOT FOUND THEN
        RAISE no_data_found;
    END IF;
    IF result.slug IS NULL
    THEN
        result.slug = '';
    END IF;
    RETURN result;
END;
$BODY$
    LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION func_get_thread_by_slug(arg_slug citext)
    RETURNS public.type_thread
AS
$BODY$
DECLARE
    result public.type_thread;
BEGIN
    SELECT *
    INTO result.id, result.slug, result.author, result.forum,
        result.title, result.message, result.votes, result.created
    FROM public.thread
    WHERE slug = arg_slug;
    result.is_new := FALSE;
    IF NOT FOUND THEN
        RAISE no_data_found;
    END IF;
    RETURN result;
END;
$BODY$
    LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION func_create_post(arg_author citext, arg_id INT, arg_message text, arg_parent INT,
                                            arg_forum citext, arg_created TIMESTAMP WITH TIME ZONE)
    RETURNS public.type_post
AS
$BODY$
DECLARE
    result          public.type_post;
    parent_thread   INT;
    author_nickname citext;
BEGIN
    SELECT nickname INTO author_nickname FROM public.person WHERE nickname = arg_author;
    IF NOT FOUND THEN
        RAISE no_data_found;
    END IF;
    SELECT thread INTO parent_thread FROM public.post WHERE id = arg_parent;
    IF arg_parent != '0' AND parent_thread != arg_id THEN
        RAISE foreign_key_violation;
    END IF;
    INSERT INTO public.post(author, thread, forum, message, parent, created)
    VALUES (arg_author, arg_id, arg_forum, arg_message, arg_parent, arg_created) RETURNING *
        INTO result.id, result.author, result.thread, result.forum,
            result.message, result.is_edited, result.parent, result.created, result.post_path;
    RETURN result;
EXCEPTION
    WHEN foreign_key_violation THEN
        RAISE foreign_key_violation;
END;
$BODY$
    LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION func_update_thread(arg_message text, arg_title text, arg_slug citext, arg_id INT)
    RETURNS public.type_thread
AS
$BODY$
DECLARE
    result public.type_thread;
BEGIN
    UPDATE public.thread
    SET message = CASE
                      WHEN arg_message != '' THEN arg_message
                      ELSE message END,
        title   = CASE
                      WHEN arg_title != '' THEN arg_title
                      ELSE title END
    WHERE slug = arg_slug
       OR id = arg_id RETURNING *
        INTO result.id, result.slug, result.author, result.forum,
            result.title, result.message, result.votes, result.created;
    result.is_new := FALSE;
    IF NOT FOUND THEN
        RAISE no_data_found;
    END IF;
    IF result.slug IS NULL
    THEN
        result.slug = '';
    END IF;
    RETURN result;
END;
$BODY$
    LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION func_create_or_update_vote(arg_user citext, arg_slug citext, arg_id INT, arg_vote INT)
    RETURNS public.type_thread
AS
$BODY$
DECLARE
    result public.type_thread;
BEGIN
    SELECT id
    INTO result.id
    FROM public.thread
    WHERE slug = arg_slug
       OR id = arg_id;
    IF NOT FOUND THEN
        RAISE no_data_found;
    END IF;

    INSERT INTO public.vote (thread_id, user_nickname, voice)
    VALUES (result.id, arg_user, arg_vote)

    ON CONFLICT ON CONSTRAINT vote_pk
        DO UPDATE
        SET voice = arg_vote
    WHERE vote.thread_id = result.id
      AND vote.user_nickname = arg_user;

    SELECT *
    INTO result.id, result.slug, result.author, result.forum,
        result.title, result.message, result.votes, result.created
    FROM public.thread
    WHERE id = result.id;
    IF NOT FOUND THEN
        RAISE no_data_found;
    END IF;
    IF result.slug IS NULL
    THEN
        result.slug = '';
    END IF;
    result.is_new := FALSE;
    RETURN result;
EXCEPTION
    WHEN foreign_key_violation THEN
        RAISE no_data_found;
END;
$BODY$
    LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION func_update_post(arg_message text, arg_id INT)
    RETURNS public.type_post
AS
$BODY$
DECLARE
    result          public.type_post;
    arg_old_message text;
BEGIN
    SELECT message INTO arg_old_message FROM public.post WHERE id = arg_id;
    UPDATE public.post
    SET message   = CASE
                        WHEN arg_message != '' THEN arg_message
                        ELSE message END,
        is_edited = CASE
                        WHEN arg_message != '' AND arg_old_message != arg_message THEN TRUE
                        ELSE FALSE END
    WHERE id = arg_id RETURNING * INTO result.id, result.author, result.thread, result.forum,
        result.message, result.is_edited, result.parent, result.created, result.post_path;
    IF NOT FOUND THEN
        RAISE no_data_found;
    END IF;
    RETURN result;
END;
$BODY$
    LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION func_get_post(arg_id INT)
    RETURNS public.type_post
AS
$BODY$
DECLARE
    result public.type_post;
BEGIN
    SELECT *
    INTO result.id, result.author, result.thread, result.forum,
        result.message, result.is_edited, result.parent, result.created, result.post_path
    FROM public.post
    WHERE id = arg_id;
    IF NOT FOUND THEN
        RAISE no_data_found;
    END IF;
    RETURN result;
END;
$BODY$
    LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION func_get_posts(arg_slug citext, arg_id INT, arg_limit INT, arg_since INT,
                                          arg_desc BOOLEAN)
    RETURNS SETOF public.type_post
AS
$BODY$
DECLARE
    result        public.type_post;
    arg_thread_id INT;
    rec           RECORD;
BEGIN
    SELECT id
    INTO arg_thread_id
    FROM public.thread
    WHERE slug = arg_slug
       OR id = arg_id;
    IF NOT FOUND THEN
        RAISE no_data_found;
    END IF;
    FOR rec IN SELECT *
               FROM public.post
               WHERE thread = arg_thread_id
                 AND CASE
                         WHEN arg_since = '0' THEN TRUE
                         WHEN arg_desc THEN id < arg_since
                         ELSE id > arg_since END
               ORDER BY (CASE WHEN arg_desc THEN created END) DESC,
                        (CASE WHEN NOT arg_desc THEN created END) ASC,
                        (CASE WHEN arg_desc THEN id END) DESC,
                        (CASE WHEN NOT arg_desc THEN id END) ASC
               LIMIT arg_limit
        LOOP
            result.id := rec.id;
            result.author := rec.author;
            result.thread := rec.thread;
            result.forum := rec.forum;
            result.message := rec.message;
            result.is_edited := rec.is_edited;
            result.parent := rec.parent;
            result.created := rec.created;
            RETURN next result;
        END LOOP;
END;
$BODY$
    LANGUAGE plpgsql;

-- CREATE OR REPLACE FUNCTION func_get_posts_flat(arg_slug citext, arg_id INT, arg_limit INT, arg_since INT,
--                                                arg_desc BOOLEAN)
--     RETURNS SETOF public.type_post
-- AS
-- $BODY$
-- DECLARE
--     result        public.type_post;
--     arg_thread_id INT;
--     rec           RECORD;
-- BEGIN
--     SELECT id
--     INTO arg_thread_id
--     FROM public.thread
--     WHERE slug = arg_slug
--        OR id = arg_id;
--     IF NOT FOUND THEN
--         RAISE no_data_found;
--     END IF;
--     FOR rec IN SELECT *
--                FROM public.post
--                WHERE thread = arg_thread_id
--                  AND CASE
--                          WHEN arg_since = '0' THEN TRUE
--                          WHEN arg_desc THEN id < arg_since
--                          ELSE id > arg_since END
--                ORDER BY (CASE WHEN arg_desc THEN id END) DESC,
--                         (CASE WHEN NOT arg_desc THEN id END) ASC
--                LIMIT arg_limit
--         LOOP
--             result.id := rec.id;
--             result.author := rec.author;
--             result.thread := rec.thread;
--             result.forum := rec.forum;
--             result.message := rec.message;
--             result.is_edited := rec.is_edited;
--             result.parent := rec.parent;
--             result.created := rec.created;
--             RETURN next result;
--         END LOOP;
-- END;
-- $BODY$
--     LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION func_get_posts_flat(arg_slug citext, arg_id INT, arg_limit INT, arg_since INT,
                                               arg_desc BOOLEAN)
    RETURNS SETOF public.type_post
AS
$BODY$
DECLARE
    result        public.type_post;
    arg_thread_id INT;
    rec           RECORD;
BEGIN
    SELECT id
    INTO arg_thread_id
    FROM public.thread
    WHERE slug = arg_slug
       OR id = arg_id;
    IF NOT FOUND THEN
        RAISE no_data_found;
    END IF;
    FOR rec IN SELECT *
               FROM public.post AS ps
                        JOIN public.person AS pr ON pr.nickname = ps.author
                        JOIN public.forum AS f ON f.slug = ps.forum
               WHERE ps.thread = arg_thread_id
                 AND CASE
                         WHEN arg_since = 0 THEN TRUE
                         ELSE CASE
                                  WHEN arg_desc THEN ps.id < arg_since
                                  ELSE ps.id > arg_since
                             END
                   END
               ORDER BY (CASE WHEN arg_desc THEN ps.id END) DESC,
                        (CASE WHEN NOT arg_desc THEN ps.id END) ASC
               LIMIT arg_limit
        LOOP
            result.id := rec.id;
            result.author := rec.author;
            result.thread := rec.thread;
            result.forum := rec.forum;
            result.message := rec.message;
            result.is_edited := rec.is_edited;
            result.parent := rec.parent;
            result.created := rec.created;
            RETURN next result;
        END LOOP;
END;
$BODY$
    LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION func_get_posts_tree(arg_slug citext, arg_id INT, arg_limit INT, arg_since INT,
                                               arg_desc BOOLEAN)
    RETURNS SETOF public.type_post
AS
$BODY$
DECLARE
    result        public.type_post;
    arg_thread_id INT;
    root_path     INT[];
    rec           RECORD;
BEGIN
    SELECT id
    INTO arg_thread_id
    FROM public.thread
    WHERE slug = arg_slug
       OR id = arg_id;
    IF NOT FOUND THEN
        RAISE no_data_found;
    END IF;
    SELECT post_path INTO root_path FROM public.post WHERE id = arg_since;
    FOR rec IN SELECT *
               FROM public.post
               WHERE thread = arg_thread_id
                 AND CASE
                         WHEN arg_since = '0' THEN TRUE
                         WHEN arg_desc THEN post_path < root_path
                         ELSE post_path > root_path END
               ORDER BY (CASE WHEN arg_desc THEN post_path END) DESC,
                        (CASE WHEN NOT arg_desc THEN post_path END) ASC
               LIMIT arg_limit
        LOOP
            result.id := rec.id;
            result.author := rec.author;
            result.thread := rec.thread;
            result.forum := rec.forum;
            result.message := rec.message;
            result.is_edited := rec.is_edited;
            result.parent := rec.parent;
            result.created := rec.created;
            result.post_path := rec.post_path;
            RETURN next result;
        END LOOP;
END;
$BODY$
    LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION func_get_posts_parent_tree(arg_slug citext, arg_id INT, arg_limit INT, arg_since INT,
                                                      arg_desc BOOLEAN)
    RETURNS SETOF public.type_post
AS
$BODY$
DECLARE
    result        public.type_post;
    arg_thread_id INT;
    rec           RECORD;
BEGIN
    SELECT id
    INTO arg_thread_id
    FROM public.thread
    WHERE slug = arg_slug
       OR id = arg_id;
    IF NOT FOUND THEN
        RAISE no_data_found;
    END IF;
    FOR rec IN
        WITH paths AS (SELECT post_path as path
                       FROM public.post
                       WHERE thread = arg_thread_id
                         AND CASE
                                 WHEN arg_since = '0' THEN TRUE
                                 WHEN arg_desc THEN post_path[2] <
                                                    (SELECT post_path[2] FROM public.post WHERE id = arg_since)
                                 ELSE post_path[2] > (SELECT post_path[2] FROM public.post WHERE id = arg_since) END)
        SELECT *
        FROM public.post
        WHERE thread = arg_thread_id
          AND post_path[2] IN (SELECT id
                               FROM public.post
                               WHERE thread = arg_thread_id
                                 AND post_path IN (SELECT path FROM paths)
                                 AND parent = 0
                               ORDER BY (CASE WHEN arg_desc THEN id END) DESC,
                                        (CASE WHEN NOT arg_desc THEN id END) ASC
                               LIMIT arg_limit)
        ORDER BY (CASE WHEN arg_desc THEN post_path[2] END) DESC, post_path,
                 (CASE WHEN NOT arg_desc THEN post_path[2] END) ASC, post_path
        LOOP
            result.id := rec.id;
            result.author := rec.author;
            result.thread := rec.thread;
            result.forum := rec.forum;
            result.message := rec.message;
            result.is_edited := rec.is_edited;
            result.parent := rec.parent;
            result.created := rec.created;
            result.post_path := rec.post_path;
            RETURN next result;
        END LOOP;
END;
$BODY$
    LANGUAGE plpgsql;

------------------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------
-- FUNCTIONS CALL ------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------

SELECT *
FROM func_add_admin();

CREATE INDEX post_rating_d_idx ON public.post USING btree (post_path DESC);

CREATE INDEX post_rating_idx ON public.post USING btree (post_path);

CREATE INDEX post_author_idx ON public.post USING btree (author);

CREATE INDEX post_thread_idx ON public.post USING btree (thread);

CREATE INDEX post_created_d_idx ON public.post USING btree (created DESC);

CREATE INDEX post_created_idx ON public.post USING btree (created);

CREATE INDEX thread_forum_idx ON public.thread USING btree (forum);

CREATE INDEX thread_author_idx ON public.thread USING btree (author);

CREATE INDEX forum_author_idx ON public.forum USING btree (author);

CREATE INDEX forum_id_idx ON public.forum USING btree (id);

CREATE INDEX forum_users_user_nickname_idx ON public.forum_users USING btree (user_nickname);

CREATE INDEX forum_users_forum_slug_idx ON public.forum_users USING btree (forum_slug);

CREATE INDEX person_id_idx ON public.forum USING btree (id);


