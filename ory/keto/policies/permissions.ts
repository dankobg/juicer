import { Namespace, SubjectSet, Context } from "@ory/keto-namespace-types";

class AdminGroup implements Namespace {
  related: {
    members: Identity[];
  };
}

class Identity implements Namespace {
  related: {
    managers: (Identity | SubjectSet<AdminGroup, "members">)[];
  };

  permits = {
    manager: (ctx: Context): boolean =>
      this.related.managers.includes(ctx.subject),
  };
}
