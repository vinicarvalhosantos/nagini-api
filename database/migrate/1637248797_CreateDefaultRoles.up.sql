-- Migration goes here.

-- Auto-generated SQL script #202111181221
INSERT INTO public.roles ("name",description,created_at,updated_at)
VALUES ('USER','The user role, can view and change information only themselves','NOW()','NOW()');


INSERT INTO public.roles ("name",description,created_at,updated_at)
VALUES ('ADMIN','The admin role, can view and change informations of anyone','NOW()','NOW()');

INSERT INTO public.roles ("name",description,created_at,updated_at)
VALUES ('SUPPORT','The support role, can view informations of anyone but can not change it','NOW()','NOW()');
