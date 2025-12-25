package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// ============================================================================
// CONFIGURATION
// ============================================================================

const (
	GraphQLEndpoint = "https://bff-page.kakao.com/graphql"

	// Thay Cookie mới của bạn vào đây
	CookieString = "_kpdid=fe3277a6e854465f872802cdf1e5891c; _kpiid=278f88bc29a85822679573a2aa34ec0f; _kahai=ba4be1e91c915262b60a6402f5e1e321446b1a1739f6a79ad21e84f1393c8927; _kau=35c5e49a69a9372b6ef6c2bc607965bf7ee6f92bf89bcc5318f18ea5e7dc7503b25394c9f33386b9c72fdd760f46731e48b12b253c105e3355f11c76a1b578a84ae02192398a0a77224fe6fc96b859b3e148c8a8dfb13c5dae2806fe5f7e51ed8478bf2de0240a01b6422618b77f387e6b3a08103120e11a21a054b33336303934353339343736373131303635343630373135333038313634363632e3a084a39eef024d31862ae3336ba294; _kawlt=KoJqnywI4FsAEHmoTpG9DLjSB9i0RC67nFInM6tZNFnium_rbz7Qm6k9wxXXgLmOdVjPOdsgThZ_K2iJtY0nJAbTsrz8ebWiTndnb7RBTZcOt2PtUYbkcIw3o8H_Kyx-; _kawltea=1766718049; _karmt=UsUxsvx5meGQq_TKLAlTO8p1ahRK_2Mjex-jM2BvpqzIA21sR2qKj7wDaWWzSHPr; _karmtea=1769245249; _kaslt=gnpr3z6X3w1YPbfcd1tZmPbLXKQDqzTLwKvruWsvTFrryq0wC921X7LrMZ6wXN3HC8hvN6oVvNhOSrK/qLChRA==; _kpawbat_e=AoN0FteLMP9ULb%2FNVd%2Bpd4oR6syQ9Y6yGdCbGplso1IBLlWGHt8YYfTVG1gHoXfV8RlCznF5mC0MJ%2BosOgj6wa1SgfGQ7YQsy356JwOFdcE%3D; _kpwtkn=eyJhbGciOiJkaXIiLCJlbmMiOiJBMjU2R0NNIn0..-ED0OT2G8NFjlrI0.S3I-JNP_eXSAUsOFyn8dAjTDyCaQv-LcHHMhUWISZePblV7fNYJIGlJvlv-PFasvIP6lRvf_Amhy-iHdSvkqbCDKkGqVRhM6eG5oTwy__RtVxpknjHyl37M00Cb8Ge1110ku2PRZbU2shFEm0DAQHjBkPvMpFnqTVripX3PHzCvrTT9wniHKpKn8Qx8AMfQULH9OoD_7a-aE2FA082-JAYp7RpstnasAFj-lkWwJWLqNqqu1M9TjHML8ktutRgWC3DJPE4-orUItTVyTSA33X_lREeUXID7l3v0APYmIKzCci8cGsGsRYz06UAAOJumVcawgmyvTkiBzeFJYH0hIt2a_XZp2oo7qusjvOD_yJ4qnAL4oYYJzqHA2TzAMx56l5ue0THFGO_hdBttADwNYnD9JkHx--g5GnzHNycP6O5IbolDkLUZkb4Woq0N2b8CLYAJyFABdVVc1eRq7T9dKGX4en7zDvNd60kTm954zljJfxU80libd3TWBYXqjG0Y14nv8qT_eUsHYD64I1BytImnJbMRxT6Yoo9hiGH3g02MA4H0zeLLG2HXzqqgNdxdUi0H9_Eow23cZulCZiRlJ4vbS3TeBi1aZ268lXQ8IBMk-62yt75Mn1xpyLW5-02iem8uYfcyIaudjhbPDTyOybV9nCdC5xZfAZ-f-poJ7GQZexyzEg0oFy9VqfQw4mzR_DF1BPB_lFtKbyShP43vcY2tOWzQzVc_1VDJ5TS86wh5GqOA8GhlsM-59si2PGzOAXICpzxQL9mUXAYWSXFJubphWHkWfZuPjFCN3dO1a0awZOtAEfFQVylQSQmKoyN-iZwYWv2Zo8f5e6pgCCjiq2zg-6yIWNA7bcijtfxTbS6sKDCj0f_QYzVWdha1mJjDo-KYU3T3wBg_dlQcccI2X8YJgDg6Jjdrko14ijxfVbLek82BU6DbJi0cJvmo0RczIr9SnJMUJjLIy0ZF-hb7Bnhp9Z60iDGH9yPJ8VAkOw3o_8pIqHXH6QDhDUz5N-y_y6dx6Va68yZ8NkpTFfApsvc_PI7gapa_ie7IDSrT-8zJEenPREfKRsRLhiVUMA2dYdPuPzWhXU-3X6vRj6x-IczkPuYGdfUdOZSHiqWcKBEbk0W3ofpG2cpdNo8uoNa8RBfiLfzxKyEVB5OVRUBfobzThYRgyjDfMsN6v98CMHxPdSHYBDABt7yyClxN0eZiBo1RR3NX6Y-psMqBDQAvpVBOJv2Z22w4X1pEH_5fQYvfhddKtv8p_6aXclmNlrq9mnJYPNDOHAMf_PfYlwxYPvUxsB5NdHo2RKgCs5Eui_qU1khn-NuKSOaNSCMLglbBbzRI0QEApAbRImsvJqhqs9hM8SV-4JrKVDEABj6nFjq9We3rVzkuJNEeBhNyT1sthUxhQCWgLPIf6urYZlAhHxzt_keLxkCHHHltxq6Nmb-VRB0E_53ZV9bMXDr_d-97TK4DI9ZklxHhyNtaPM64rVCh0wOEfoc9jbD1zTO9WuGnPF4DOmNKFd-T4tj4f8FswLimAEBCUl7ykQjGPKu8-lpy5hkXGHeXVFlaNJ_U9FUqb7T4hh4F5AoO5Ng4_u1EkkdMgMQF0aFMVn9lJe1D9_Mb4xJ-mTGAm1eGXS-OKvfduX97kGSvOVbw_Q3ecEE8SqjACvJGnBPwTjpijCTAjGQuiebJW8mI7o8HS_vpV3iVP9cuJs2K8RR-JWBhPj95B3NdAc16Z812Tmksmeoac0xevbfiKUeA5moiXybZSpMl6KYbu-GnvHMw-VX-uI5g3MDV9iVPb3xFBvlmzIjnWrv7_WuQRVoze335CPbXNem2SKuDXZylNZ-86yre3MQ-CC3poto-OZS_JR3qnWCOcTU5-X_1Ityvo-jKVOHbFmkVf7UpGBeRot6dlNY8SLPwj5xZoYkScoCKOTQqJ45Vr7vxe1NgwZmb5XIz8O-xd1Kcfteei4u6SjRBRvwm5Zb1VNHFcS9OGlIa01RD9Vupf25oR_KquYqFRgdXfR33sWuPjKUkPD46dq-j5GqdcCFUMZK9NVkce9WG_Uq46-0So3pp8rYgvDvn_TZ5jY48q55nb85T8a-wsiCviBr43we8kw-CaDvrQgmM9hDm7HOPj2orjPQUb35fIhQel_SoNzjoPtmc-h1Rk1O3wzAFgrojbWEZ2pPv2JlIQe_sFhTg6QO1EssObHKcerPmO8H7Y2uVuBSHywVLTKu1uNW9FIshiHHq435R0P2KvtNGfKRj_ekBFtcbFnZhkpBmJ34lvPnl4JNRVWZZeB9ig_1tI3H8JpNqeTLA.XaFrZ2WqiUcmrsRXqQricw; _T_ANO=EbYxru9UaTh9x2uRBgT6H7LXSvoba6tDNb0ZTXkvsToiGNJD9j59X56rSwnF1x/obmLyxG/bknGlpc9OoOgSRPt9Zt6KDQ0Uocd0CTa2caAf6jzdrEe5ECuiTScBKEj8gH3TyhJg8Nh4IgTNuwzqDk1I+3WTpVkBq+wwsablZrHQeYEpMf5WL/YfxZtUCNYwOneLVHpuelR+xQXSFeLxCJ+APr7hxdDsGVWusWkA2vj+tG7ozggh5ov1iw+r1woAPSzm9pVtc7M0gogBvIHWkv6fG46bwFgCha8IFoAybXNWIl6KMMsh+rxoSMSL/UAoJ82bmynMOxGo/GO5rot0hQ=="
)

// ============================================================================
// QUERY (dùng payload bạn cung cấp — giữ nguyên các type Long!)
// ============================================================================

const ViewerQuery = `
query viewerInfo($seriesId: Long!, $productId: Long!) {
  viewerInfo(seriesId: $seriesId, productId: $productId) {
    item {
      ...SingleFragment
    }
    seriesItem {
      ...SeriesFragment
    }
    prevItem {
      ...NearItemFragment
    }
    nextItem {
      ...NearItemFragment
    }
    viewerData {
      ...TextViewerData
      ...TalkViewerData
      ...ImageViewerData
      ...VodViewerData
    }
    displayAd {
      ...DisplayAd
    }
  }
}

fragment SingleFragment on Single {
  id
  productId
  seriesId
  title
  thumbnail
  badge
  isFree
  ageGrade
  state
  slideType
  lastReleasedDate
  size
  pageCount
  isHidden
  remainText
  isWaitfreeBlocked
  saleState
  series {
    ...SeriesFragment
  }
  serviceProperty {
    ...ServicePropertyFragment
  }
  operatorProperty {
    ...OperatorPropertyFragment
  }
  assetProperty {
    ...AssetPropertyFragment
  }
  discountRate
  discountRateText
  isShortsDrama
}

fragment SeriesFragment on Series {
  id
  seriesId
  title
  thumbnail
  landThumbnail
  categoryUid
  lang
  category
  categoryType
  subcategoryUid
  subcategory
  badge
  isAllFree
  isWaitfree
  ageGrade
  state
  onIssue
  authors
  description
  pubPeriod
  freeSlideCount
  lastSlideAddedDate
  waitfreeBlockCount
  waitfreePeriodByMinute
  bm
  saleState
  startSaleDt
  saleMethod
  discountRate
  discountRateText
  serviceProperty {
    ...ServicePropertyFragment
  }
  operatorProperty {
    ...OperatorPropertyFragment
  }
  assetProperty {
    ...AssetPropertyFragment
  }
  translateProperty {
    ...TranslatePropertyFragment
  }
}

fragment ServicePropertyFragment on ServiceProperty {
  viewCount
  readCount
  ratingCount
  ratingSum
  commentCount
  pageContinue {
    ...ContinueInfoFragment
  }
  todayGift {
    ...TodayGift
  }
  preview {
    ...PreviewFragment
  }
  waitfreeTicket {
    ...WaitfreeTicketFragment
  }
  isAlarmOn
  isLikeOn
  ticketCount
  purchasedDate
  lastViewInfo {
    ...LastViewInfoFragment
  }
  purchaseInfo {
    ...PurchaseInfoFragment
  }
  ticketInfo {
    price
    discountPrice
    ticketType
  }
}

fragment ContinueInfoFragment on ContinueInfo {
  title
  isFree
  productId
  lastReadProductId
  scheme
  continueProductType
  hasNewSingle
  hasUnreadSingle
}

fragment TodayGift on TodayGift {
  id
  uid
  ticketType
  ticketKind
  ticketCount
  ticketExpireAt
  ticketExpiredText
  isReceived
  seriesId
}

fragment PreviewFragment on Preview {
  item {
    ...PreviewSingleFragment
  }
  nextItem {
    ...PreviewSingleFragment
  }
  usingScroll
}

fragment PreviewSingleFragment on Single {
  id
  productId
  seriesId
  title
  thumbnail
  badge
  isFree
  ageGrade
  state
  slideType
  lastReleasedDate
  size
  pageCount
  isHidden
  remainText
  isWaitfreeBlocked
  saleState
  operatorProperty {
    ...OperatorPropertyFragment
  }
  assetProperty {
    ...AssetPropertyFragment
  }
}

fragment OperatorPropertyFragment on OperatorProperty {
  thumbnail
  copy
  helixImpId
  isTextViewer
  selfCensorship
  isBook
  cashInfo {
    discountRate
    setDiscountRate
  }
  ticketInfo {
    price
    discountPrice
    ticketType
  }
}

fragment AssetPropertyFragment on AssetProperty {
  bannerImage
  cardImage
  cardTextImage
  cleanImage
  ipxVideo
  bannerSet {
    ...BannerSetFragment
  }
  cardSet {
    ...CardSetFragment
  }
  cardCover {
    ...CardCoverFragment
  }
}

fragment BannerSetFragment on BannerSet {
  backgroundImage
  backgroundColor
  mainImage
  titleImage
}

fragment CardSetFragment on CardSet {
  backgroundColor
  backgroundImage
}

fragment CardCoverFragment on CardCover {
  coverImg
  coverRestricted
}

fragment WaitfreeTicketFragment on WaitfreeTicket {
  chargedPeriod
  chargedCount
  chargedAt
}

fragment LastViewInfoFragment on LastViewInfo {
  isDone
  lastViewDate
  rate
  spineIndex
}

fragment PurchaseInfoFragment on PurchaseInfo {
  purchaseType
  rentExpireDate
  expired
}

fragment TranslatePropertyFragment on TranslateProperty {
  category {
    ...LocaleMapFragment
  }
  sub_category {
    ...LocaleMapFragment
  }
}

fragment LocaleMapFragment on LocaleMap {
  ko
  en
  th
}

fragment NearItemFragment on NearItem {
  productId
  slideType
  ageGrade
  isFree
  title
  thumbnail
}

fragment TextViewerData on TextViewerData {
  type
  atsServerUrl
  metaSecureUrl
  contentsList {
    chapterId
    contentId
    secureUrl
  }
}

fragment TalkViewerData on TalkViewerData {
  type
  talkDownloadData {
    dec
    host
    path
    talkViewerType
  }
}

fragment ImageViewerData on ImageViewerData {
  type
  imageDownloadData {
    ...ImageDownloadData
  }
}

fragment ImageDownloadData on ImageDownloadData {
  files {
    ...ImageDownloadFile
  }
  totalCount
  totalSize
  viewDirection
  gapBetweenImages
  readType
}

fragment ImageDownloadFile on ImageDownloadFile {
  no
  size
  secureUrl
  width
  height
}

fragment VodViewerData on VodViewerData {
  type
  vodDownloadData {
    contentPackId
    drmType
    endpointUrl
    width
    height
    duration
  }
  drmInfo {
    type
    serverType
    error
    fairplayCertificateUrl
    widevineLicenseUrl
    fairplayLicenseUrl
    token
    provider
    assertion
  }
}

fragment DisplayAd on DisplayAd {
  sectionUid
  bannerUid
  treviUid
  momentUid
}
`

// ============================================================================
// MODELS (mình chỉ khai báo những field cần dùng)
// ============================================================================

type GraphQLPayload struct {
	OperationName string                 `json:"operationName"`
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
}

type GQLResponse struct {
	Data struct {
		ViewerInfo ViewerInfo `json:"viewerInfo"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type ViewerInfo struct {
	Item       Item       `json:"item"`
	SeriesItem *Series    `json:"seriesItem,omitempty"`
	PrevItem   *NearItem  `json:"prevItem,omitempty"`
	NextItem   *NearItem  `json:"nextItem,omitempty"`
	ViewerData ViewerData `json:"viewerData"`
}

type Item struct {
	Title string `json:"title"`
	// thêm field nếu cần về sau
	ProductId int `json:"productId,omitempty"`
	SeriesId  int `json:"seriesId,omitempty"`
}

type Series struct {
	SeriesId int    `json:"seriesId"`
	Title    string `json:"title"`
}

type NearItem struct {
	ProductId int    `json:"productId"`
	Title     string `json:"title"`
	// slideType, etc. omitted — thêm nếu cần
}

type ViewerData struct {
	AtsServerUrl string        `json:"atsServerUrl"`
	MetaSecure   string        `json:"metaSecureUrl,omitempty"`
	ContentsList []ContentPart `json:"contentsList"`
}

type ContentPart struct {
	ChapterId int    `json:"chapterId,omitempty"`
	ContentId int    `json:"contentId,omitempty"`
	SecureUrl string `json:"secureUrl"`
}

type ContentFileResponse struct {
	ContentInfo struct {
		ParagraphList []ParagraphNode `json:"paragraphList"`
	} `json:"contentInfo"`
}

type ParagraphNode struct {
	Type     string          `json:"type"`
	Text     string          `json:"text"`
	Children []ParagraphNode `json:"childParagraphList,omitempty"`
}

// ============================================================================
// MAIN: loop theo nextItem.productId — giữ flow tải/chạy y nguyên code gốc
// ============================================================================

func main() {
	// context tổng cho toàn bộ chạy
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	client := &http.Client{Timeout: 30 * time.Second}
	headers := map[string]string{
		"accept":          "*/*",
		"accept-language": "vi-VN,vi;q=0.9,en-US;q=0.8,en;q=0.7",
		"content-type":    "application/json",
		"origin":          "https://page.kakao.com",
		"referer":         "https://page.kakao.com/",
		"user-agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36",
		"cookie":          CookieString,
	}

	// Bắt đầu từ đây: set seriesId và productId chương ban đầu bạn muốn crawl
	currentSeries := 61097897
	currentProduct := 66155147

	for {
		fmt.Printf("\n>>> Crawl productId = %d (seriesId=%d)\n", currentProduct, currentSeries)

		// 1) build payload giữ nguyên như bạn yêu cầu
		payload := GraphQLPayload{
			OperationName: "viewerInfo",
			Query:         ViewerQuery,
			Variables: map[string]interface{}{
				"seriesId":  currentSeries,
				"productId": currentProduct,
			},
		}
		payloadBytes, _ := json.Marshal(payload)

		respBody, err := makeRequest(ctx, client, "POST", GraphQLEndpoint, headers, bytes.NewBuffer(payloadBytes))
		if err != nil {
			fmt.Printf(">>> ERROR: GraphQL request failed: %v\n", err)
			break
		}

		var gqlResp GQLResponse
		if err := json.Unmarshal(respBody, &gqlResp); err != nil {
			fmt.Printf(">>> ERROR: Parse Error: %v\nBody: %s\n", err, string(respBody))
			break
		}
		if len(gqlResp.Errors) > 0 {
			fmt.Printf(">>> ERROR: GraphQL Error: %s\n", gqlResp.Errors[0].Message)
			break
		}

		viewerData := gqlResp.Data.ViewerInfo.ViewerData
		itemInfo := gqlResp.Data.ViewerInfo.Item
		nextItem := gqlResp.Data.ViewerInfo.NextItem

		if viewerData.AtsServerUrl == "" {
			fmt.Println(">>> ERROR: Không tìm thấy atsServerUrl. Kiểm tra lại Cookie.")
			break
		}

		// Tạo file
		safeTitle := sanitizeFilename(itemInfo.Title)
		fileName := safeTitle + ".txt"

		file, err := os.Create(fileName)
		if err != nil {
			fmt.Printf(">>> ERROR: Không thể tạo file %s: %v\n", fileName, err)
			break
		}
		fmt.Printf(">>> Đã tạo file: %s\n", fileName)
		fmt.Printf(">>> Tổng số phần: %d. Bắt đầu tải...\n", len(viewerData.ContentsList))

		// Ghi tiêu đề
		if _, err := file.WriteString(fmt.Sprintf("# %s\n\n", itemInfo.Title)); err != nil {
			fmt.Printf(">>> WARNING: lỗi ghi tiêu đề: %v\n", err)
		}

		// Loop tải các phần
		for i, part := range viewerData.ContentsList {
			if i > 0 {
				time.Sleep(1 * time.Second)
			}
			fullUrl := viewerData.AtsServerUrl + part.SecureUrl
			fmt.Printf(" -> [%d/%d] Đang tải... (Delay 1s)\n", i+1, len(viewerData.ContentsList))

			partBody, err := makeRequest(ctx, client, "GET", fullUrl, headers, nil)
			if err != nil {
				fmt.Printf("   [Lỗi tải phần %d]: %v\n", i+1, err)
				continue
			}

			var contentResp ContentFileResponse
			if err := json.Unmarshal(partBody, &contentResp); err != nil {
				// nếu không parse được JSON -> có thể là image / binary -> bỏ qua
				continue
			}

			for _, node := range contentResp.ContentInfo.ParagraphList {
				text := extractTextRecursive(node)
				text = strings.ReplaceAll(text, "&nbsp;", " ")
				if _, err := file.WriteString(text); err != nil {
					fmt.Printf("   Lỗi ghi file: %v\n", err)
				}
			}
		}

		file.Close()
		fmt.Println(">>> XONG! Kiểm tra file:", fileName)

		// Nếu không có nextItem -> dừng
		if nextItem == nil || nextItem.ProductId == 0 {
			fmt.Println(">>> HẾT CHƯƠNG hoặc không tìm thấy nextItem.")
			break
		}

		// Cập nhật productId để crawl chương tiếp theo
		fmt.Printf(">>> Next: %d (%s)\n", nextItem.ProductId, nextItem.Title)
		currentProduct = nextItem.ProductId

		// delay an toàn trước lần tiếp theo
		time.Sleep(2 * time.Second)
	}
}

// ============================================================================
// HELPERS
// ============================================================================

func makeRequest(ctx context.Context, client *http.Client, method, url string, headers map[string]string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Status %d | Body: %s", resp.StatusCode, string(respBytes))
	}
	return respBytes, nil
}

func extractTextRecursive(node ParagraphNode) string {
	var sb strings.Builder
	if node.Text != "" {
		sb.WriteString(node.Text)
	}
	if node.Type == "BR" {
		sb.WriteString("\n")
	}
	for _, child := range node.Children {
		sb.WriteString(extractTextRecursive(child))
	}
	switch node.Type {
	case "P", "DIV":
		sb.WriteString("\n")
	case "H3":
		sb.WriteString("\n\n### ")
	}
	return sb.String()
}

func sanitizeFilename(name string) string {
	reg, _ := regexp.Compile(`[\\/:*?"<>|]`)
	safe := reg.ReplaceAllString(name, "_")
	return strings.TrimSpace(safe)
}
