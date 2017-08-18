package cmd

import (
	"fmt"
	"os"
	"path"
	"github.com/Unknwon/com"
	"github.com/go-macaron/binding"
	"github.com/urfave/cli"
	"gopkg.in/macaron.v1"

	"github.com/gityflow/githorse/models"
	"github.com/gityflow/githorse/pkg/context"
	"github.com/gityflow/githorse/pkg/form"
	"github.com/gityflow/githorse/pkg/setting"
	"github.com/gityflow/githorse/routes"
	"github.com/gityflow/githorse/routes/admin"
	apiv1 "github.com/gityflow/githorse/routes/api/v1"
	"github.com/gityflow/githorse/routes/dev"
	"github.com/gityflow/githorse/routes/org"
	"github.com/gityflow/githorse/routes/repo"
	"github.com/gityflow/githorse/routes/user"
)

func routesV1(c *cli.Context,m *macaron.Macaron, reqSignIn, ignSignIn, ignSignInAndCsrf, reqSignOut macaron.Handler) {
	bindIgnErr := binding.BindIgnErr
	// FIXME: not all routes need go through same middlewares.
	// Especially some AJAX requests, we can reduce middleware number to improve performance.
	// Routers.
	m.Get("/", ignSignIn, routes.Home)
	m.Group("/explore", func() {
		m.Get("", func(c *context.Context) {
			c.Redirect(setting.AppSubURL + "/explore/repos")
		})
		m.Get("/repos", routes.ExploreRepos)
		m.Get("/users", routes.ExploreUsers)
		m.Get("/organizations", routes.ExploreOrganizations)
	}, ignSignIn)
	m.Combo("/install", routes.InstallInit).Get(routes.Install).
		Post(bindIgnErr(form.Install{}), routes.InstallPost)
	m.Get("/^:type(issues|pulls)$", reqSignIn, user.Issues)

	// ***** START: User *****
	m.Group("/user", func() {
		m.Group("/login", func() {
			m.Combo("").Get(user.Login).
				Post(bindIgnErr(form.SignIn{}), user.LoginPost)
			m.Combo("/two_factor").Get(user.LoginTwoFactor).Post(user.LoginTwoFactorPost)
			m.Combo("/two_factor_recovery_code").Get(user.LoginTwoFactorRecoveryCode).Post(user.LoginTwoFactorRecoveryCodePost)
		})

		m.Get("/sign_up", user.SignUp)
		m.Post("/sign_up", bindIgnErr(form.Register{}), user.SignUpPost)
		m.Get("/reset_password", user.ResetPasswd)
		m.Post("/reset_password", user.ResetPasswdPost)
	}, reqSignOut)

	m.Group("/user/settings", func() {
		m.Get("", user.Settings)
		m.Post("", bindIgnErr(form.UpdateProfile{}), user.SettingsPost)
		m.Combo("/avatar").Get(user.SettingsAvatar).
			Post(binding.MultipartForm(form.Avatar{}), user.SettingsAvatarPost)
		m.Post("/avatar/delete", user.SettingsDeleteAvatar)
		m.Combo("/email").Get(user.SettingsEmails).
			Post(bindIgnErr(form.AddEmail{}), user.SettingsEmailPost)
		m.Post("/email/delete", user.DeleteEmail)
		m.Get("/password", user.SettingsPassword)
		m.Post("/password", bindIgnErr(form.ChangePassword{}), user.SettingsPasswordPost)
		m.Combo("/ssh").Get(user.SettingsSSHKeys).
			Post(bindIgnErr(form.AddSSHKey{}), user.SettingsSSHKeysPost)
		m.Post("/ssh/delete", user.DeleteSSHKey)
		m.Group("/security", func() {
			m.Get("", user.SettingsSecurity)
			m.Combo("/two_factor_enable").Get(user.SettingsTwoFactorEnable).
				Post(user.SettingsTwoFactorEnablePost)
			m.Combo("/two_factor_recovery_codes").Get(user.SettingsTwoFactorRecoveryCodes).
				Post(user.SettingsTwoFactorRecoveryCodesPost)
			m.Post("/two_factor_disable", user.SettingsTwoFactorDisable)
		})
		m.Group("/repositories", func() {
			m.Get("", user.SettingsRepos)
			m.Post("/leave", user.SettingsLeaveRepo)
		})
		m.Group("/organizations", func() {
			m.Get("", user.SettingsOrganizations)
			m.Post("/leave", user.SettingsLeaveOrganization)
		})
		m.Combo("/applications").Get(user.SettingsApplications).
			Post(bindIgnErr(form.NewAccessToken{}), user.SettingsApplicationsPost)
		m.Post("/applications/delete", user.SettingsDeleteApplication)
		m.Route("/delete", "GET,POST", user.SettingsDelete)
	}, reqSignIn, func(c *context.Context) {
		c.Data["PageIsUserSettings"] = true
	})

	m.Group("/user", func() {
		m.Any("/activate", user.Activate)
		m.Any("/activate_email", user.ActivateEmail)
		m.Get("/email2user", user.Email2User)
		m.Get("/forget_password", user.ForgotPasswd)
		m.Post("/forget_password", user.ForgotPasswdPost)
		m.Get("/logout", user.SignOut)
	})
	// ***** END: User *****

	adminReq := context.Toggle(&context.ToggleOptions{SignInRequired: true, AdminRequired: true})

	// ***** START: Admin *****
	m.Group("/admin", func() {
		m.Get("", adminReq, admin.Dashboard)
		m.Get("/config", admin.Config)
		m.Post("/config/test_mail", admin.SendTestMail)
		m.Get("/monitor", admin.Monitor)

		m.Group("/users", func() {
			m.Get("", admin.Users)
			m.Combo("/new").Get(admin.NewUser).Post(bindIgnErr(form.AdminCrateUser{}), admin.NewUserPost)
			m.Combo("/:userid").Get(admin.EditUser).Post(bindIgnErr(form.AdminEditUser{}), admin.EditUserPost)
			m.Post("/:userid/delete", admin.DeleteUser)
		})

		m.Group("/orgs", func() {
			m.Get("", admin.Organizations)
		})

		m.Group("/repos", func() {
			m.Get("", admin.Repos)
			m.Post("/delete", admin.DeleteRepo)
		})

		m.Group("/auths", func() {
			m.Get("", admin.Authentications)
			m.Combo("/new").Get(admin.NewAuthSource).Post(bindIgnErr(form.Authentication{}), admin.NewAuthSourcePost)
			m.Combo("/:authid").Get(admin.EditAuthSource).
				Post(bindIgnErr(form.Authentication{}), admin.EditAuthSourcePost)
			m.Post("/:authid/delete", admin.DeleteAuthSource)
		})

		m.Group("/notices", func() {
			m.Get("", admin.Notices)
			m.Post("/delete", admin.DeleteNotices)
			m.Get("/empty", admin.EmptyNotices)
		})
	}, adminReq)
	// ***** END: Admin *****

	m.Group("", func() {
		m.Group("/:username", func() {
			m.Get("", user.Profile)
			m.Get("/followers", user.Followers)
			m.Get("/following", user.Following)
			m.Get("/stars", user.Stars)
		})

		m.Get("/attachments/:uuid", func(c *context.Context) {
			attach, err := models.GetAttachmentByUUID(c.Params(":uuid"))
			if err != nil {
				c.NotFoundOrServerError("GetAttachmentByUUID", models.IsErrAttachmentNotExist, err)
				return
			} else if !com.IsFile(attach.LocalPath()) {
				c.NotFound()
				return
			}

			fr, err := os.Open(attach.LocalPath())
			if err != nil {
				c.Handle(500, "Open", err)
				return
			}
			defer fr.Close()

			c.Header().Set("Cache-Control", "public,max-age=86400")
			fmt.Println("attach.Name:", attach.Name)
			c.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, attach.Name))
			if err = repo.ServeData(c, attach.Name, fr); err != nil {
				c.Handle(500, "ServeData", err)
				return
			}
		})
		m.Post("/issues/attachments", repo.UploadIssueAttachment)
		m.Post("/releases/attachments", repo.UploadReleaseAttachment)
	}, ignSignIn)

	m.Group("/:username", func() {
		m.Get("/action/:action", user.Action)
	}, reqSignIn)

	if macaron.Env == macaron.DEV {
		m.Get("/template/*", dev.TemplatePreview)
	}

	reqRepoAdmin := context.RequireRepoAdmin()
	reqRepoWriter := context.RequireRepoWriter()

	// ***** START: Organization *****
	m.Group("/org", func() {
		m.Group("", func() {
			m.Get("/create", org.Create)
			m.Post("/create", bindIgnErr(form.CreateOrg{}), org.CreatePost)
		}, func(c *context.Context) {
			if !c.User.CanCreateOrganization() {
				c.NotFound()
			}
		})

		m.Group("/:org", func() {
			m.Get("/dashboard", user.Dashboard)
			m.Get("/^:type(issues|pulls)$", user.Issues)
			m.Get("/members", org.Members)
			m.Get("/members/action/:action", org.MembersAction)

			m.Get("/teams", org.Teams)
		}, context.OrgAssignment(true))

		m.Group("/:org", func() {
			m.Get("/teams/:team", org.TeamMembers)
			m.Get("/teams/:team/repositories", org.TeamRepositories)
			m.Route("/teams/:team/action/:action", "GET,POST", org.TeamsAction)
			m.Route("/teams/:team/action/repo/:action", "GET,POST", org.TeamsRepoAction)
		}, context.OrgAssignment(true, false, true))

		m.Group("/:org", func() {
			m.Get("/teams/new", org.NewTeam)
			m.Post("/teams/new", bindIgnErr(form.CreateTeam{}), org.NewTeamPost)
			m.Get("/teams/:team/edit", org.EditTeam)
			m.Post("/teams/:team/edit", bindIgnErr(form.CreateTeam{}), org.EditTeamPost)
			m.Post("/teams/:team/delete", org.DeleteTeam)

			m.Group("/settings", func() {
				m.Combo("").Get(org.Settings).
					Post(bindIgnErr(form.UpdateOrgSetting{}), org.SettingsPost)
				m.Post("/avatar", binding.MultipartForm(form.Avatar{}), org.SettingsAvatar)
				m.Post("/avatar/delete", org.SettingsDeleteAvatar)

				m.Group("/hooks", func() {
					m.Get("", org.Webhooks)
					m.Post("/delete", org.DeleteWebhook)
					m.Get("/:type/new", repo.WebhooksNew)
					m.Post("/gogs/new", bindIgnErr(form.NewWebhook{}), repo.WebHooksNewPost)
					m.Post("/slack/new", bindIgnErr(form.NewSlackHook{}), repo.SlackHooksNewPost)
					m.Post("/discord/new", bindIgnErr(form.NewDiscordHook{}), repo.DiscordHooksNewPost)
					m.Get("/:id", repo.WebHooksEdit)
					m.Post("/gogs/:id", bindIgnErr(form.NewWebhook{}), repo.WebHooksEditPost)
					m.Post("/slack/:id", bindIgnErr(form.NewSlackHook{}), repo.SlackHooksEditPost)
					m.Post("/discord/:id", bindIgnErr(form.NewDiscordHook{}), repo.DiscordHooksEditPost)
				})

				m.Route("/delete", "GET,POST", org.SettingsDelete)
			})

			m.Route("/invitations/new", "GET,POST", org.Invitation)
		}, context.OrgAssignment(true, true))
	}, reqSignIn)
	// ***** END: Organization *****

	// ***** START: Repository *****
	m.Group("/repo", func() {
		m.Get("/create", repo.Create)
		m.Post("/create", bindIgnErr(form.CreateRepo{}), repo.CreatePost)
		m.Get("/migrate", repo.Migrate)
		m.Post("/migrate", bindIgnErr(form.MigrateRepo{}), repo.MigratePost)
		m.Combo("/fork/:repoid").Get(repo.Fork).
			Post(bindIgnErr(form.CreateRepo{}), repo.ForkPost)
	}, reqSignIn)

	m.Group("/:username/:reponame", func() {
		m.Group("/settings", func() {
			m.Combo("").Get(repo.Settings).
				Post(bindIgnErr(form.RepoSetting{}), repo.SettingsPost)
			m.Group("/collaboration", func() {
				m.Combo("").Get(repo.SettingsCollaboration).Post(repo.SettingsCollaborationPost)
				m.Post("/access_mode", repo.ChangeCollaborationAccessMode)
				m.Post("/delete", repo.DeleteCollaboration)
			})
			m.Group("/branches", func() {
				m.Get("", repo.SettingsBranches)
				m.Post("/default_branch", repo.UpdateDefaultBranch)
				m.Combo("/*").Get(repo.SettingsProtectedBranch).
					Post(bindIgnErr(form.ProtectBranch{}), repo.SettingsProtectedBranchPost)
			}, func(c *context.Context) {
				if c.Repo.Repository.IsMirror {
					c.NotFound()
					return
				}
			})

			m.Group("/hooks", func() {
				m.Get("", repo.Webhooks)
				m.Post("/delete", repo.DeleteWebhook)
				m.Get("/:type/new", repo.WebhooksNew)
				m.Post("/gogs/new", bindIgnErr(form.NewWebhook{}), repo.WebHooksNewPost)
				m.Post("/slack/new", bindIgnErr(form.NewSlackHook{}), repo.SlackHooksNewPost)
				m.Post("/discord/new", bindIgnErr(form.NewDiscordHook{}), repo.DiscordHooksNewPost)
				m.Post("/gogs/:id", bindIgnErr(form.NewWebhook{}), repo.WebHooksEditPost)
				m.Post("/slack/:id", bindIgnErr(form.NewSlackHook{}), repo.SlackHooksEditPost)
				m.Post("/discord/:id", bindIgnErr(form.NewDiscordHook{}), repo.DiscordHooksEditPost)

				m.Group("/:id", func() {
					m.Get("", repo.WebHooksEdit)
					m.Post("/test", repo.TestWebhook)
					m.Post("/redelivery", repo.RedeliveryWebhook)
				})

				m.Group("/git", func() {
					m.Get("", repo.SettingsGitHooks)
					m.Combo("/:name").Get(repo.SettingsGitHooksEdit).
						Post(repo.SettingsGitHooksEditPost)
				}, context.GitHookService())
			})

			m.Group("/keys", func() {
				m.Combo("").Get(repo.SettingsDeployKeys).
					Post(bindIgnErr(form.AddSSHKey{}), repo.SettingsDeployKeysPost)
				m.Post("/delete", repo.DeleteDeployKey)
			})

		}, func(c *context.Context) {
			c.Data["PageIsSettings"] = true
		})
	}, reqSignIn, context.RepoAssignment(), reqRepoAdmin, context.RepoRef())

	m.Get("/:username/:reponame/action/:action", reqSignIn, context.RepoAssignment(), repo.Action)
	m.Group("/:username/:reponame", func() {
		m.Get("/issues", repo.RetrieveLabels, repo.Issues)
		m.Get("/issues/:index", repo.ViewIssue)
		m.Get("/labels/", repo.RetrieveLabels, repo.Labels)
		m.Get("/milestones", repo.Milestones)
	}, ignSignIn, context.RepoAssignment(true))
	m.Group("/:username/:reponame", func() {
		// FIXME: should use different URLs but mostly same logic for comments of issue and pull reuqest.
		// So they can apply their own enable/disable logic on routers.
		m.Group("/issues", func() {
			m.Combo("/new", repo.MustEnableIssues).Get(context.RepoRef(), repo.NewIssue).
				Post(bindIgnErr(form.NewIssue{}), repo.NewIssuePost)

			m.Group("/:index", func() {
				m.Post("/title", repo.UpdateIssueTitle)
				m.Post("/content", repo.UpdateIssueContent)
				m.Combo("/comments").Post(bindIgnErr(form.CreateComment{}), repo.NewComment)
			})
		})
		m.Group("/comments/:id", func() {
			m.Post("", repo.UpdateCommentContent)
			m.Post("/delete", repo.DeleteComment)
		})
	}, reqSignIn, context.RepoAssignment(true))
	m.Group("/:username/:reponame", func() {
		m.Group("/wiki", func() {
			m.Get("/?:page", repo.Wiki)
			m.Get("/_pages", repo.WikiPages)
		}, repo.MustEnableWiki, context.RepoRef())
	}, ignSignIn, context.RepoAssignment(false, true))

	m.Group("/:username/:reponame", func() {
		// FIXME: should use different URLs but mostly same logic for comments of issue and pull reuqest.
		// So they can apply their own enable/disable logic on routers.
		m.Group("/issues", func() {
			m.Group("/:index", func() {
				m.Post("/label", repo.UpdateIssueLabel)
				m.Post("/milestone", repo.UpdateIssueMilestone)
				m.Post("/assignee", repo.UpdateIssueAssignee)
			}, reqRepoWriter)
		})
		m.Group("/labels", func() {
			m.Post("/new", bindIgnErr(form.CreateLabel{}), repo.NewLabel)
			m.Post("/edit", bindIgnErr(form.CreateLabel{}), repo.UpdateLabel)
			m.Post("/delete", repo.DeleteLabel)
			m.Post("/initialize", bindIgnErr(form.InitializeLabels{}), repo.InitializeLabels)
		}, reqRepoWriter, context.RepoRef())
		m.Group("/milestones", func() {
			m.Combo("/new").Get(repo.NewMilestone).
				Post(bindIgnErr(form.CreateMilestone{}), repo.NewMilestonePost)
			m.Get("/:id/edit", repo.EditMilestone)
			m.Post("/:id/edit", bindIgnErr(form.CreateMilestone{}), repo.EditMilestonePost)
			m.Get("/:id/:action", repo.ChangeMilestonStatus)
			m.Post("/delete", repo.DeleteMilestone)
		}, reqRepoWriter, context.RepoRef())

		m.Group("/releases", func() {
			m.Get("/new", repo.NewRelease)
			m.Post("/new", bindIgnErr(form.NewRelease{}), repo.NewReleasePost)
			m.Post("/delete", repo.DeleteRelease)
			m.Get("/edit/*", repo.EditRelease)
			m.Post("/edit/*", bindIgnErr(form.EditRelease{}), repo.EditReleasePost)
		}, repo.MustBeNotBare, reqRepoWriter, func(c *context.Context) {
			c.Data["PageIsViewFiles"] = true
		})

		// FIXME: Should use c.Repo.PullRequest to unify template, currently we have inconsistent URL
		// for PR in same repository. After select branch on the page, the URL contains redundant head user name.
		// e.g. /org1/test-repo/compare/master...org1:develop
		// which should be /org1/test-repo/compare/master...develop
		m.Combo("/compare/*", repo.MustAllowPulls).Get(repo.CompareAndPullRequest).
			Post(bindIgnErr(form.NewIssue{}), repo.CompareAndPullRequestPost)

		m.Group("", func() {
			m.Combo("/_edit/*").Get(repo.EditFile).
				Post(bindIgnErr(form.EditRepoFile{}), repo.EditFilePost)
			m.Combo("/_new/*").Get(repo.NewFile).
				Post(bindIgnErr(form.EditRepoFile{}), repo.NewFilePost)
			m.Post("/_preview/*", bindIgnErr(form.EditPreviewDiff{}), repo.DiffPreviewPost)
			m.Combo("/_delete/*").Get(repo.DeleteFile).
				Post(bindIgnErr(form.DeleteRepoFile{}), repo.DeleteFilePost)

			m.Group("", func() {
				m.Combo("/_upload/*").Get(repo.UploadFile).
					Post(bindIgnErr(form.UploadRepoFile{}), repo.UploadFilePost)
				m.Post("/upload-file", repo.UploadFileToServer)
				m.Post("/upload-remove", bindIgnErr(form.RemoveUploadFile{}), repo.RemoveUploadFileFromServer)
			}, func(c *context.Context) {
				if !setting.Repository.Upload.Enabled {
					c.NotFound()
					return
				}
			})
		}, repo.MustBeNotBare, reqRepoWriter, context.RepoRef(), func(c *context.Context) {
			if !c.Repo.CanEnableEditor() {
				c.NotFound()
				return
			}

			c.Data["PageIsViewFiles"] = true
		})
	}, reqSignIn, context.RepoAssignment())

	m.Group("/:username/:reponame", func() {
		m.Group("", func() {
			m.Get("/releases", repo.MustBeNotBare, repo.Releases)
			m.Get("/pulls", repo.RetrieveLabels, repo.Pulls)
			m.Get("/pulls/:index", repo.ViewPull)
		}, context.RepoRef())

		m.Group("/branches", func() {
			m.Get("", repo.Branches)
			m.Get("/all", repo.AllBranches)
			m.Post("/delete/*", reqSignIn, reqRepoWriter, repo.DeleteBranchPost)
			m.Post("/create/*", reqSignIn, reqRepoWriter, repo.CreateNewBranch)
		}, repo.MustBeNotBare, func(c *context.Context) {
			c.Data["PageIsViewFiles"] = true
		})

		m.Group("/wiki", func() {
			m.Group("", func() {
				m.Combo("/_new").Get(repo.NewWiki).
					Post(bindIgnErr(form.NewWiki{}), repo.NewWikiPost)
				m.Combo("/:page/_edit").Get(repo.EditWiki).
					Post(bindIgnErr(form.NewWiki{}), repo.EditWikiPost)
				m.Post("/:page/delete", repo.DeleteWikiPagePost)
			}, reqSignIn, reqRepoWriter)
		}, repo.MustEnableWiki, context.RepoRef())

		m.Get("/archive/*", repo.MustBeNotBare, repo.Download)

		m.Group("/pulls/:index", func() {
			m.Get("/commits", context.RepoRef(), repo.ViewPullCommits)
			m.Get("/files", context.RepoRef(), repo.ViewPullFiles)
			m.Post("/merge", reqRepoWriter, repo.MergePullRequest)
		}, repo.MustAllowPulls)

		m.Group("", func() {
			m.Get("/src/*", repo.Home)
			m.Get("/raw/*", repo.SingleDownload)
			m.Get("/commits/*", repo.RefCommits)
			m.Get("/commit/:sha([a-f0-9]{7,40})$", repo.Diff)
			m.Get("/forks", repo.Forks)
		}, repo.MustBeNotBare, context.RepoRef())
		m.Get("/commit/:sha([a-f0-9]{7,40})\\.:ext(patch|diff)", repo.MustBeNotBare, repo.RawDiff)

		m.Get("/compare/:before([a-z0-9]{40})\\.\\.\\.:after([a-z0-9]{40})", repo.MustBeNotBare, context.RepoRef(), repo.CompareDiff)
	}, ignSignIn, context.RepoAssignment())
	m.Group("/:username/:reponame", func() {
		m.Get("/stars", repo.Stars)
		m.Get("/watchers", repo.Watchers)
	}, ignSignIn, context.RepoAssignment(), context.RepoRef())

	m.Group("/:username", func() {
		m.Get("/:reponame", ignSignIn, context.RepoAssignment(), context.RepoRef(), repo.Home)

		m.Group("/:reponame", func() {
			m.Head("/tasks/trigger", repo.TriggerTask)
		})
		// Use the regexp to match the repository name
		// Duplicated routes to enable different ways of accessing same set of URLs,
		// e.g. with or without ".git" suffix.
		m.Group("/:reponame([\\d\\w-_\\.]+\\.git$)", func() {
			m.Get("", ignSignIn, context.RepoAssignment(), context.RepoRef(), repo.Home)
			m.Route("/*", "GET,POST", ignSignInAndCsrf, repo.HTTPContexter(), repo.HTTP)
		})
		m.Route("/:reponame/*", "GET,POST", ignSignInAndCsrf, repo.HTTPContexter(), repo.HTTP)
	})
	// ***** END: Repository *****

	m.Group("/api", func() {
		apiv1.RegisterRoutes(m)
	}, ignSignIn)

	// robots.txt
	m.Get("/robots.txt", func(c *context.Context) {
		if setting.HasRobotsTxt {
			c.ServeFileContent(path.Join(setting.CustomPath, "robots.txt"))
		} else {
			c.NotFound()
		}
	})

}
